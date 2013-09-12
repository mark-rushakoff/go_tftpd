package session_router

import (
	"testing"

	"github.com/mark-rushakoff/go_tftpd/readsessioncollection"
	"github.com/mark-rushakoff/go_tftpd/safepackets"
	"github.com/mark-rushakoff/go_tftpd/safety_filter"
	"github.com/mark-rushakoff/go_tftpd/test_helpers"
	"github.com/mark-rushakoff/go_tftpd/timeout_controller"
)

func TestRouteAckRoutes(t *testing.T) {
	sessions := readsessioncollection.NewReadSessionCollection()
	router := NewSessionRouter(sessions)
	fakeAddr := test_helpers.MakeMockAddr("fake_network", "a")

	acks := make(chan *safepackets.SafeAck, 1)
	timeoutController := &timeout_controller.MockTimeoutController{
		HandleAckHandler: func(ack *safepackets.SafeAck) {
			acks <- ack
		},
	}
	sessions.Add(timeoutController, fakeAddr)

	router.RouteAck(&safety_filter.IncomingSafeAck{
		Addr: fakeAddr,
		Ack:  safepackets.NewSafeAck(8),
	})

	select {
	case ack := <-acks:
		if ack.BlockNumber != 8 {
			t.Fatalf("Received incorrect ack")
		}
	default:
		t.Fatalf("RouteAck should have sent Ack")
	}
}

func TestRouteAckToMissingSessionDoesNotPanic(t *testing.T) {
	sessions := readsessioncollection.NewReadSessionCollection()
	router := NewSessionRouter(sessions)
	fakeAddr := test_helpers.MakeMockAddr("fake_network", "a")

	router.RouteAck(&safety_filter.IncomingSafeAck{
		Addr: fakeAddr,
		Ack:  safepackets.NewSafeAck(8),
	})

	// ok
}
