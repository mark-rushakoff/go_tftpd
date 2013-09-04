package session_router

import (
	"testing"

	"github.com/mark-rushakoff/go_tftpd/read_session_collection"
	"github.com/mark-rushakoff/go_tftpd/safe_packets"
	"github.com/mark-rushakoff/go_tftpd/safety_filter"
	"github.com/mark-rushakoff/go_tftpd/test_helpers"
	"github.com/mark-rushakoff/go_tftpd/timeout_controller"
)

func TestRouteAckRoutes(t *testing.T) {
	sessions := read_session_collection.NewReadSessionCollection()
	router := NewSessionRouter(sessions)
	fakeAddr := test_helpers.MakeMockAddr("fake_network", "a")

	acks := make(chan *safe_packets.SafeAck, 1)
	timeoutController := &timeout_controller.MockTimeoutController{
		HandleAckHandler: func(ack *safe_packets.SafeAck) {
			acks <- ack
		},
	}
	sessions.Add(timeoutController, fakeAddr)

	router.RouteAck(&safety_filter.IncomingSafeAck{
		Addr: fakeAddr,
		Ack:  safe_packets.NewSafeAck(8),
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
