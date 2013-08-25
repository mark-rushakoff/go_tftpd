package request_router

import (
	"testing"
	"time"

	"github.com/mark-rushakoff/go_tftpd/safe_packets"
	"github.com/mark-rushakoff/go_tftpd/safety_filter"
	"github.com/mark-rushakoff/go_tftpd/session_manager"
	"github.com/mark-rushakoff/go_tftpd/test_helpers"
)

func TestReadsRouteToCreateSession(t *testing.T) {
	sessionRead := make(chan *safety_filter.IncomingSafeReadRequest)
	sessionManager := &session_manager.SessionManager{
		ReadRequest: sessionRead,
	}

	safetyFilter := makeSafetyFilter()

	router := NewRequestRouter(safetyFilter, sessionManager)
	go router.Route()

	fakeIncoming := &safety_filter.IncomingSafeReadRequest{
		Read: &safe_packets.SafeReadRequest{},
		Addr: test_helpers.MakeMockAddr("fake_net", "a"),
	}

	go func() {
		safetyFilter.IncomingRead <- fakeIncoming
	}()

	select {
	case routedRead := <-sessionRead:
		if routedRead != fakeIncoming {
			t.Errorf("Expected RequestRouter to route incoming read request, but it did not")
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
