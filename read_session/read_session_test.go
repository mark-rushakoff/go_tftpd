package read_session

import (
	"bytes"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/mark-rushakoff/go_tftpd/response_agent"
	"github.com/mark-rushakoff/go_tftpd/safe_packets"
	"github.com/mark-rushakoff/go_tftpd/timeout_controller"
)

type suite struct {
	responseAgent     *response_agent.MockResponseAgent
	timeoutController *timeout_controller.MockTimeoutController
	config            *ReadSessionConfig
	session           *ReadSession
}

func makeSuite(blockSize uint16, readerContent string) *suite {
	timeoutController := timeout_controller.MakeMockTimeoutController()
	responseAgent := response_agent.MakeMockResponseAgent()
	config := &ReadSessionConfig{
		ResponseAgent:     responseAgent,
		Reader:            strings.NewReader(readerContent),
		BlockSize:         blockSize,
		TimeoutController: timeoutController,
	}
	return &suite{
		timeoutController: timeoutController,
		responseAgent:     responseAgent,
		config:            config,
		session:           NewReadSession(config),
	}
}

func TestBeginStartsLifecycle(t *testing.T) {
	s := makeSuite(3, "foobar")

	assertCountdownCalls(t, s.timeoutController, 0)
	go s.session.Begin()

	runtime.Gosched() // TODO: should have a better way to synchronize here

	assertCountdownCalls(t, s.timeoutController, 1)
	assertTotalMessagesSent(t, s.responseAgent, 1)

	sentData := s.responseAgent.MostRecentData()
	assertDataMessage(t, sentData, 1, []byte("foo"))
}

func TestAckLastDataCausesFinish(t *testing.T) {
	s := makeSuite(64, "Hello!")

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.session.Begin()
	}()
	s.session.Ack() <- safe_packets.NewSafeAck(1)

	select {
	case <-s.session.Finished():
	// ok
	case <-time.After(5 * time.Millisecond):
		t.Fatalf("Session did not finish upon receiving terminal ack")
	}

	terminated := make(chan bool)
	go func() {
		wg.Wait()
		terminated <- true
	}()

	select {
	case <-terminated:
	// ok
	case <-time.After(time.Millisecond):
		t.Fatalf("Finishing did not terminate the goroutine")
	}
}

func TestAckMultipleDataCausesFinish(t *testing.T) {
	s := makeSuite(5, "foobar")

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.session.Begin()
	}()
	s.session.Ack() <- safe_packets.NewSafeAck(1)
	assertCountdownCalls(t, s.timeoutController, 2) // one for first data sent, one for data sent after ack

	assertTotalMessagesSent(t, s.responseAgent, 2)

	sentData := s.responseAgent.MostRecentData()
	assertDataMessage(t, sentData, 2, []byte("r"))
	assertRestartCalls(t, s.timeoutController, 1)

	s.session.Ack() <- safe_packets.NewSafeAck(2)
	assertCountdownCalls(t, s.timeoutController, 2) // no need to adjust countdown calls; was ack to last packet

	select {
	case <-s.session.Finished():
	// ok
	case <-time.After(5 * time.Millisecond):
		t.Fatalf("Session did not finish upon receiving terminal ack")
	}

	terminated := make(chan bool)
	go func() {
		wg.Wait()
		terminated <- true
	}()

	select {
	case <-terminated:
	// ok
	case <-time.After(time.Millisecond):
		t.Fatalf("Finishing did not terminate the goroutine")
	}
}

func TestOldAckCausesDataResend(t *testing.T) {
	s := makeSuite(2, "foobar")

	go s.session.Begin()
	runtime.Gosched()

	s.session.Ack() <- safe_packets.NewSafeAck(1)
	s.responseAgent.Reset()
	runtime.Gosched() // TODO: find better way to synchronize

	assertCountdownCalls(t, s.timeoutController, 2) // one for first data sent, one for data sent after ack
	assertTotalMessagesSent(t, s.responseAgent, 1)
	sentData := s.responseAgent.MostRecentData()
	assertDataMessage(t, sentData, 2, []byte("ob"))
	assertRestartCalls(t, s.timeoutController, 1)

	s.session.Ack() <- safe_packets.NewSafeAck(1)
	s.responseAgent.Reset()
	runtime.Gosched() // TODO: find better way to synchronize

	assertCountdownCalls(t, s.timeoutController, 3) // +1 for re-sending old data
	assertRestartCalls(t, s.timeoutController, 2)   // +1 for re-sending old data
	assertTotalMessagesSent(t, s.responseAgent, 1)

	sentData = s.responseAgent.MostRecentData()
	assertDataMessage(t, sentData, 2, []byte("ob"))
}

func TestTimeoutCausesDataResend(t *testing.T) {
	s := makeSuite(3, "foobar")

	go s.session.Begin()

	runtime.Gosched() // TODO: find better way to synchronize

	// first data
	assertTotalMessagesSent(t, s.responseAgent, 1)
	assertCountdownCalls(t, s.timeoutController, 1) // one for first data sent

	sentData := s.responseAgent.MostRecentData()
	assertDataMessage(t, sentData, 1, []byte("foo"))

	s.responseAgent.Reset()

	s.timeoutController.CauseNonExpiredTimeout()

	runtime.Gosched() // TODO: find better way to synchronize

	// repeated data
	assertCountdownCalls(t, s.timeoutController, 2) // +1 for re-sending data
	assertTotalMessagesSent(t, s.responseAgent, 1)

	sentData = s.responseAgent.MostRecentData()
	assertDataMessage(t, sentData, 1, []byte("foo"))
}

func TestExpiredTimeoutCausesFinished(t *testing.T) {
	s := makeSuite(64, "Hello!")

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.session.Begin()
	}()

	runtime.Gosched() // TODO: find better way to synchronize

	s.responseAgent.Reset()
	s.timeoutController.CauseExpiredTimeout()

	select {
	case <-s.session.Finished():
	// ok
	case <-time.After(5 * time.Millisecond):
		t.Fatalf("Session did not finish upon receiving expired timeout")
	}

	terminated := make(chan bool)
	go func() {
		wg.Wait()
		terminated <- true
	}()

	select {
	case <-terminated:
	// ok
	case <-time.After(time.Millisecond):
		t.Fatalf("Finishing by expired timeout did not terminate the goroutine")
	}

	assertTotalMessagesSent(t, s.responseAgent, 0)
}

/*
func TestMultipleDataPackets(t *testing.T) {
	responseAgent := response_agent.MakeMockResponseAgent()

	config := &ReadSessionConfig{
		ResponseAgent:     responseAgent,
		Reader:            strings.NewReader("12345678abcdef"),
		BlockSize:         8,
		TimeoutController: timeout_controller.MakeMockTimeoutController(),
	}
	session := NewReadSession(config)
	session.Begin()

	assertTotalMessagesSent(t, responseAgent, 2)

	sentData := responseAgent.MostRecentData()
	assertDataMessage(t, sentData, 1, []byte("12345678"))

	responseAgent.Reset()
	session.Ack <- safe_packets.NewSafeAck(1)

	sentData = responseAgent.MostRecentData()
	assertDataMessage(t, sentData, 2, []byte("abcdef"))

	session.assertNotFinished(t)
}

*/

func assertDataMessage(t *testing.T, data *safe_packets.SafeData, expectedBlockNumber uint16, expectedData []byte) {
	if data == nil {
		_, file, line, _ := runtime.Caller(1)
		t.Fatalf("Data not sent at %v:%v", file, line)
	}

	actualBlockNumber := data.BlockNumber
	if actualBlockNumber != expectedBlockNumber {
		_, file, line, _ := runtime.Caller(1)
		t.Errorf("Expected ReadSession to send data with block number %v, received %v at %v:%v", expectedBlockNumber, actualBlockNumber, file, line)
	}

	if !bytes.Equal(data.Data.Data, expectedData) {
		_, file, line, _ := runtime.Caller(1)
		t.Errorf("Expected ReadSession to send data %v, received %v at %v:%v", expectedData, data.Data.Data, file, line)
	}
}

func assertTotalMessagesSent(t *testing.T, responseAgent *response_agent.MockResponseAgent, total int) {
	actualTotal := responseAgent.TotalMessagesSent()
	if actualTotal != total {
		_, file, line, _ := runtime.Caller(1)
		t.Fatalf("Expected %v message(s) sent but %v message(s) were sent at %v:%v", total, actualTotal, file, line)
	}
}

func assertCountdownCalls(t *testing.T, timeoutController *timeout_controller.MockTimeoutController, numCalls uint) {
	actualCount := timeoutController.CountdownCalls()
	if actualCount != numCalls {
		_, file, line, _ := runtime.Caller(1)
		t.Fatalf("Expected %v call(s) to Countdown but %v call(s) received at %v:%v", numCalls, actualCount, file, line)
	}
}

func assertRestartCalls(t *testing.T, timeoutController *timeout_controller.MockTimeoutController, numCalls uint) {
	actualCount := timeoutController.RestartCalls()
	if actualCount != numCalls {
		_, file, line, _ := runtime.Caller(1)
		t.Fatalf("Expected %v call(s) to Restart but %v call(s) received at %v:%v", numCalls, actualCount, file, line)
	}
}
