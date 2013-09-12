package sessionrouter

import (
	"testing"

	"github.com/mark-rushakoff/go_tftpd/readsessioncollection"
	"github.com/mark-rushakoff/go_tftpd/safepackets"
	"github.com/mark-rushakoff/go_tftpd/safetyfilter"
	"github.com/mark-rushakoff/go_tftpd/testhelpers"
	"github.com/mark-rushakoff/go_tftpd/timeoutcontroller"
)

func TestRouteAckRoutes(t *testing.T) {
	sessions := readsessioncollection.NewReadSessionCollection()
	router := NewSessionRouter(sessions)
	fakeAddr := testhelpers.MakeMockAddr("fake_network", "a")

	acks := make(chan *safepackets.SafeAck, 1)
	timeoutController := &timeoutcontroller.MockTimeoutController{
		HandleAckHandler: func(ack *safepackets.SafeAck) {
			acks <- ack
		},
	}
	sessions.Add(timeoutController, fakeAddr)

	router.RouteAck(&safetyfilter.IncomingSafeAck{
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
	fakeAddr := testhelpers.MakeMockAddr("fake_network", "a")

	router.RouteAck(&safetyfilter.IncomingSafeAck{
		Addr: fakeAddr,
		Ack:  safepackets.NewSafeAck(8),
	})

	// ok
}
