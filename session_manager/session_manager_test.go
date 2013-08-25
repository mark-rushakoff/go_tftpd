package session_manager

import (
	"testing"
	"time"

	"github.com/mark-rushakoff/go_tftpd/read_session"
	"github.com/mark-rushakoff/go_tftpd/safe_packets"
	"github.com/mark-rushakoff/go_tftpd/safety_filter"
	"github.com/mark-rushakoff/go_tftpd/test_helpers"
)

func TestIncomingReadCreatesSession(t *testing.T) {
	factoryCalls := make(chan bool, 1)
	readSessionFactory := func(*read_session.ReadSessionConfig) *read_session.ReadSession {
		factoryCalls <- true
		return &read_session.ReadSession{}
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
