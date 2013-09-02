package timeout_controller

import (
	"testing"
	"time"

	"github.com/mark-rushakoff/go_tftpd/read_session"
	"github.com/mark-rushakoff/go_tftpd/safe_packets"
)

func TestBeginThenTimeoutResendsData(t *testing.T) {
	begin := make(chan bool, 1)
	resend := make(chan bool, 1)
	session := &read_session.MockReadSession{
		BeginHandler: func() {
			begin <- true
		},
		ResendHandler: func() {
			resend <- true
		},
	}
	controller := NewTimeoutController(3*time.Millisecond, 3, session)
	controller.Begin()

	select {
	case <-begin:
	// ok
	case <-time.After(time.Millisecond):
		t.Fatalf("Controller did not call session.Begin")
	}

	select {
	case <-time.After(2 * time.Millisecond):
		// ok
	case <-resend:
		t.Fatalf("Controller called resend too early")
	}

	select {
	case <-resend:
		// ok
	case <-time.After(time.Millisecond):
		t.Fatalf("Controller did not call resend in time")
	}
}

func TestNewDoesNotStartTimer(t *testing.T) {
	resend := make(chan bool, 1)
	session := &read_session.MockReadSession{
		ResendHandler: func() {
			resend <- true
		},
	}
	NewTimeoutController(3*time.Millisecond, 3, session)

	select {
	case <-time.After(4 * time.Millisecond):
		// ok
	case <-resend:
		t.Fatalf("Controller re-sent before begin called")
	}
}

func TestAckResetsTimer(t *testing.T) {
	resend := make(chan bool, 1)
	ack := make(chan *safe_packets.SafeAck, 1)
	session := &read_session.MockReadSession{
		BeginHandler: func() {
		},
		ResendHandler: func() {
			resend <- true
		},
		HandleAckHandler: func(a *safe_packets.SafeAck) {
			ack <- a
		},
	}
	controller := NewTimeoutController(3*time.Millisecond, 3, session)
	controller.Begin()

	time.Sleep(2 * time.Millisecond)
	controller.HandleAck(safe_packets.NewSafeAck(8))

	select {
	case a := <-ack:
		if a.BlockNumber != 8 {
			t.Errorf("Controller sent ack with wrong block number")
		}
	default:
		t.Fatalf("Controller did not forward ack to session")
	}

	select {
	case <-time.After(2 * time.Millisecond):
		// ok
	case <-resend:
		t.Fatalf("Controller re-sent data too early after timer was supposed to be reset")
	}

	select {
	case <-resend:
		// ok
	case <-time.After(time.Millisecond):
		t.Fatalf("Controller did not re-send data on time")
	}
}

func TestStopResendingAfterTryLimit(t *testing.T) {
	resend := make(chan bool, 1)
	session := &read_session.MockReadSession{
		BeginHandler: func() {
		},
		ResendHandler: func() {
			resend <- true
		},
	}
	controller := NewTimeoutController(3*time.Millisecond, 3, session)
	controller.Begin() // begin contains first try

	select {
	case <-resend:
		// ok, second try
	case <-time.After(4 * time.Millisecond):
		t.Fatalf("Controller did not re-send in time")
	}

	select {
	case <-resend:
		// ok, third try
	case <-time.After(4 * time.Millisecond):
		t.Fatalf("Controller did not re-send in time")
	}

	select {
	case <-time.After(4 * time.Millisecond):
		// ok, correct timeout
	case <-resend:
		t.Fatalf("Controller re-sent when tries should have been exhausted")
	}
}

func TestHandleAckResetsTryLimit(t *testing.T) {
	resend := make(chan bool, 1)
	session := &read_session.MockReadSession{
		BeginHandler: func() {
		},
		ResendHandler: func() {
			resend <- true
		},
		HandleAckHandler: func(_ *safe_packets.SafeAck) {
		},
	}
	controller := NewTimeoutController(3*time.Millisecond, 3, session)
	controller.Begin() // begin contains first try

	select {
	case <-resend:
		// ok, second try
	case <-time.After(4 * time.Millisecond):
		t.Fatalf("Controller did not re-send in time")
	}

	select {
	case <-resend:
		// ok, third try
	case <-time.After(4 * time.Millisecond):
		t.Fatalf("Controller did not re-send in time")
	}

	controller.HandleAck(safe_packets.NewSafeAck(5))
	for i := 0; i < 3; i++ {
		select {
		case <-resend:
			// ok, try 1-2-3
		case <-time.After(4 * time.Millisecond):
			t.Fatalf("Controller did not re-send in time")
		}
	}

	select {
	case <-time.After(4 * time.Millisecond):
		// ok, correct timeout
	case <-resend:
		t.Fatalf("Controller re-sent when tries should have been exhausted")
	}
}
