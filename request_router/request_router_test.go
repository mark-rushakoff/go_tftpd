package request_router

import (
	"testing"
	"time"

	"github.com/mark-rushakoff/go_tftpd/safe_packets"
	"github.com/mark-rushakoff/go_tftpd/safety_filter"
	"github.com/mark-rushakoff/go_tftpd/session_manager"
	"github.com/mark-rushakoff/go_tftpd/test_helpers"
)

func TestRoutesProperlyToSession(t *testing.T) {
	sessionAck := make(chan *safety_filter.IncomingSafeAck)
	sessionRead := make(chan *safety_filter.IncomingSafeReadRequest)
	sessionManager := &session_manager.SessionManager{
		ReadRequest: sessionRead,
		Ack:         sessionAck,
	}

	safetyFilter := makeSafetyFilter()

	router := NewRequestRouter(safetyFilter, sessionManager)
	go router.Route()

	fakeIncomingRead := &safety_filter.IncomingSafeReadRequest{
		Read: &safe_packets.SafeReadRequest{},
		Addr: test_helpers.MakeMockAddr("fake_net", "a"),
	}

	go func() {
		safetyFilter.IncomingRead <- fakeIncomingRead
	}()

	select {
	case routedRead := <-sessionRead:
		if routedRead != fakeIncomingRead {
			t.Errorf("Expected RequestRouter to route incoming read request, but it did not")
		}
	case <-time.After(10 * time.Millisecond):
		t.Errorf("Request not routed in time")
	}

	fakeIncomingAck := &safety_filter.IncomingSafeAck{
		Ack:  &safe_packets.SafeAck{},
		Addr: test_helpers.MakeMockAddr("fake_net", "a"),
	}

	go func() {
		safetyFilter.IncomingAck <- fakeIncomingAck
	}()

	select {
	case routedAck := <-sessionAck:
		if routedAck != fakeIncomingAck {
			t.Errorf("Expected RequestRouter to route incoming ack, but it did not")
		}
	case <-time.After(10 * time.Millisecond):
		t.Errorf("Request not routed in time")
	}
}

func makeSafetyFilter() *safety_filter.SafetyFilter {
	return &safety_filter.SafetyFilter{
		IncomingAck:  make(chan *safety_filter.IncomingSafeAck),
		IncomingRead: make(chan *safety_filter.IncomingSafeReadRequest),
	}
}
