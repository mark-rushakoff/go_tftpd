package read_session

import (
	"bytes"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/mark-rushakoff/go_tftpd/response_agent"
	"github.com/mark-rushakoff/go_tftpd/safe_packets"
)

func TestBegin(t *testing.T) {
	responseAgent := response_agent.MakeMockResponseAgent()

	config := ReadSessionConfig{
		ResponseAgent: responseAgent,
		Reader:        strings.NewReader("Hello!"),
		BlockSize:     512,
	}
	session := NewReadSession(config)
	session.Begin()

	assertTotalMessagesSent(t, responseAgent, 2)

	sentAck := responseAgent.MostRecentAck()
	actualBlockNumber := sentAck.BlockNumber
	expectedBlockNumber := uint16(0)
	if actualBlockNumber != expectedBlockNumber {
		t.Errorf("Expected ReadSession to ack with block number %v, received %v", expectedBlockNumber, actualBlockNumber)
		bytes.Equal(nil, nil)
	}

	sentData := responseAgent.MostRecentData()
	assertDataMessage(t, sentData, 1, []byte("Hello!"))
}

func TestMultipleDataPackets(t *testing.T) {
	responseAgent := response_agent.MakeMockResponseAgent()

	config := ReadSessionConfig{
		ResponseAgent: responseAgent,
		Reader:        strings.NewReader("12345678abcdef"),
		BlockSize:     8,
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
}

func TestOldAck(t *testing.T) {
	responseAgent := response_agent.MakeMockResponseAgent()

	config := ReadSessionConfig{
		ResponseAgent: responseAgent,
		Reader:        strings.NewReader("12345678abcdefgh876543210"),
		BlockSize:     8,
	}
	session := NewReadSession(config)
	session.Begin()

	assertTotalMessagesSent(t, responseAgent, 2)

	sentData := responseAgent.MostRecentData()
	assertDataMessage(t, sentData, 1, []byte("12345678"))

	responseAgent.Reset()
	session.Ack <- safe_packets.NewSafeAck(1)

	sentData = responseAgent.MostRecentData()
	assertDataMessage(t, sentData, 2, []byte("abcdefgh"))

	assertTotalMessagesSent(t, responseAgent, 1)

	responseAgent.Reset()
	session.Ack <- safe_packets.NewSafeAck(1)

	// yield to the session's channel... probably a better way to do this? Or maybe it's just a test artifact?
	time.Sleep(1 * time.Millisecond)

	assertTotalMessagesSent(t, responseAgent, 1)

	sentData = responseAgent.MostRecentData()
	assertDataMessage(t, sentData, 2, []byte("abcdefgh"))
}

func assertDataMessage(t *testing.T, data *safe_packets.SafeData, expectedBlockNumber uint16, expectedData []byte) {
	if data == nil {
		t.Fatalf("Data not sent")
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
