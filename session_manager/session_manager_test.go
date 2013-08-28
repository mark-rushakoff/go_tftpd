package session_manager

import (
	"net"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/mark-rushakoff/go_tftpd/read_session"
	"github.com/mark-rushakoff/go_tftpd/response_agent"
	"github.com/mark-rushakoff/go_tftpd/safe_packets"
	"github.com/mark-rushakoff/go_tftpd/safety_filter"
	"github.com/mark-rushakoff/go_tftpd/test_helpers"
	"github.com/mark-rushakoff/go_tftpd/timeout_controller"
)

func TestIncomingReadCreatesSession(t *testing.T) {
	factoryCalls := make(chan bool, 1)
	readSessionFactory := func(filename string, _ net.Addr) *read_session.ReadSession {
		factoryCalls <- true
		return &read_session.ReadSession{
			Config: &read_session.ReadSessionConfig{
				ResponseAgent:     response_agent.MakeMockResponseAgent(),
				Reader:            strings.NewReader(filename),
				BlockSize:         1,
				TimeoutController: timeout_controller.MakeMockTimeoutController(),
			},
		}
	}
	manager := NewSessionManager(readSessionFactory)
	go manager.Watch()

	read := safe_packets.NewSafeReadRequest("foobar", safe_packets.NetAscii)
	addr := test_helpers.MakeMockAddr("fakenet", "a")

	manager.ReadRequest <- &safety_filter.IncomingSafeReadRequest{
		Read: read,
		Addr: addr,
	}

	select {
	case <-factoryCalls:
		// ok
	case <-time.After(5 * time.Millisecond):
		t.Errorf("Did not see read session created in time")
	}
}

func TestIncomingAckIsRoutedToCorrectSession(t *testing.T) {
	responseAgentA := response_agent.MakeMockResponseAgent()
	readSessionA := &read_session.ReadSession{
		Config: &read_session.ReadSessionConfig{
			ResponseAgent:     responseAgentA,
			Reader:            strings.NewReader("Abc"),
			BlockSize:         1,
			TimeoutController: timeout_controller.MakeMockTimeoutController(),
		},
		Ack: make(chan *safe_packets.SafeAck),
	}

	responseAgentB := response_agent.MakeMockResponseAgent()
	readSessionB := &read_session.ReadSession{
		Config: &read_session.ReadSessionConfig{
			ResponseAgent:     responseAgentB,
			Reader:            strings.NewReader("Def"),
			BlockSize:         1,
			TimeoutController: timeout_controller.MakeMockTimeoutController(),
		},
		Ack: make(chan *safe_packets.SafeAck),
	}
	totalCalls := 0
	readSessionFactory := func(filename string, _ net.Addr) *read_session.ReadSession {
		totalCalls++
		if totalCalls == 1 {
			return readSessionA
		} else if totalCalls == 2 {
			return readSessionB
		} else {
			panic("test Read session factory called too many times")
		}
	}
	manager := NewSessionManager(readSessionFactory)
	go manager.Watch()

	manager.ReadRequest <- &safety_filter.IncomingSafeReadRequest{
		Read: safe_packets.NewSafeReadRequest("foobar", safe_packets.NetAscii),
		Addr: test_helpers.MakeMockAddr("fakenet", "a"),
	}

	expectedData := safe_packets.NewSafeData(1, []byte("A"))
	if !responseAgentA.MostRecentData().Equals(expectedData) {
		t.Fatalf("Received incorrect data %v, expected %v", responseAgentA.MostRecentData(), expectedData)
	}

	manager.ReadRequest <- &safety_filter.IncomingSafeReadRequest{
		Read: safe_packets.NewSafeReadRequest("baz", safe_packets.NetAscii),
		Addr: test_helpers.MakeMockAddr("fakenet", "b"),
	}

	runtime.Gosched()

	expectedData = safe_packets.NewSafeData(1, []byte("D"))
	if !responseAgentB.MostRecentData().Equals(expectedData) {
		t.Fatalf("Received incorrect data %v, expected %v", responseAgentB.MostRecentData(), expectedData)
	}

	responseAgentA.Reset()
	responseAgentB.Reset()

	go func() {
		manager.Ack <- &safety_filter.IncomingSafeAck{
			Ack:  safe_packets.NewSafeAck(1),
			Addr: test_helpers.MakeMockAddr("fakenet", "a"),
		}
	}()

	runtime.Gosched()

	expectedData = safe_packets.NewSafeData(2, []byte("b"))
	if !responseAgentA.MostRecentData().Equals(expectedData) {
		t.Fatalf("Received incorrect data %v, expected %v", responseAgentA.MostRecentData(), expectedData)
	}

	go func() {
		manager.Ack <- &safety_filter.IncomingSafeAck{
			Ack:  safe_packets.NewSafeAck(1),
			Addr: test_helpers.MakeMockAddr("fakenet", "b"),
		}
	}()

	runtime.Gosched()

	expectedData = safe_packets.NewSafeData(2, []byte("e"))
	if !responseAgentB.MostRecentData().Equals(expectedData) {
		t.Fatalf("Received incorrect data %v, expected %v", responseAgentB.MostRecentData(), expectedData)
	}
}
