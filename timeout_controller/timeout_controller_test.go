package timeout_controller

import (
	"runtime"
	"testing"

	"github.com/mark-rushakoff/go_tftpd/readsession"
	"github.com/mark-rushakoff/go_tftpd/safepackets"
)

func TestBeginSessionThenTimeoutResendsData(t *testing.T) {
	begin := make(chan bool, 1)
	resend := make(chan bool, 1)
	session := &readsession.MockReadSession{
		BeginHandler: func() {
			begin <- true
		},
		ResendHandler: func() {
			resend <- true
		},
	}
	timer := NewMockTimer(make(chan bool, 1), make(chan bool, 1))
	controller := manualTimeoutController(3, session, func() {}, timer)
	controller.BeginSession()

	select {
	case <-begin:
	// ok
	default:
		t.Fatalf("Controller did not call session.BeginSession")
	}

	select {
	case <-resend:
		t.Fatalf("Controller called resend too early")
	default:
		// ok
	}

	timer.Elapse()

	select {
	case <-resend:
		// ok
	default:
		t.Fatalf("Controller did not call resend when timer elapsed")
	}
}

func TestNewDoesNotStartTimer(t *testing.T) {
	resend := make(chan bool, 1)
	restartTimer := make(chan bool, 1)
	session := &readsession.MockReadSession{
		ResendHandler: func() {
			resend <- true
		},
	}
	timer := NewMockTimer(restartTimer, make(chan bool, 1))
	manualTimeoutController(3, session, func() {}, timer)

	select {
	case <-resend:
		t.Fatalf("Controller re-sent before begin called")
	case <-restartTimer:
		t.Fatalf("Timer started before begin called")
	default:
		// ok
	}
}

func TestAckRestartsTimer(t *testing.T) {
	resend := make(chan bool, 1)
	restartTimer := make(chan bool, 1)
	ack := make(chan *safepackets.SafeAck, 1)
	session := &readsession.MockReadSession{
		BeginHandler: func() {
		},
		ResendHandler: func() {
			resend <- true
		},
		HandleAckHandler: func(a *safepackets.SafeAck) {
			ack <- a
		},
	}
	timer := NewMockTimer(restartTimer, make(chan bool, 1))
	controller := manualTimeoutController(3, session, func() {}, timer)
	controller.BeginSession()
	select {
	case <-restartTimer:
		// ok
	default:
		t.Fatalf("Timer should have been restarted upon begin")
	}

	controller.HandleAck(safepackets.NewSafeAck(8))
	select {
	case <-restartTimer:
		// ok
	default:
		t.Fatalf("Timer should have been restarted upon handling ack")
	}
	select {
	case a := <-ack:
		if a.BlockNumber != 8 {
			t.Errorf("Controller sent ack with wrong block number")
		}
	default:
		t.Fatalf("Controller did not forward ack to session")
	}

	timer.Elapse()

	select {
	case <-resend:
		// ok
	default:
		t.Fatalf("Controller did not re-send data after timer elapsed")
	}
}

func TestStopResendingAfterTryLimit(t *testing.T) {
	send := make(chan bool, 1)
	restartTimer := make(chan bool, 1)
	session := &readsession.MockReadSession{
		BeginHandler: func() {
			send <- true
		},
		ResendHandler: func() {
			send <- true
		},
	}
	timer := NewMockTimer(restartTimer, make(chan bool, 1))
	controller := manualTimeoutController(3, session, func() {}, timer)
	// 3 tries remaining
	controller.BeginSession()
	select {
	case <-restartTimer:
		// ok
	default:
		t.Fatalf("Timer should have been restarted upon begin")
	}
	select {
	case <-send:
	// ok
	default:
		t.Fatalf("Controller should have sent data upon begin")
	}

	// 2 tries remaining
	timer.Elapse()
	select {
	case <-restartTimer:
		// ok
	default:
		t.Fatalf("Timer should have been restarted upon first elapse")
	}
	select {
	case <-send:
		// ok
	default:
		t.Fatalf("Controller should have re-sent after first elapse")
	}

	// last try
	timer.Elapse()
	runtime.Gosched()
	select {
	case <-send:
		// ok, third try
	default:
		t.Fatalf("Controller should have re-sent after second elapse")
	}
	select {
	case <-restartTimer:
		// ok
	default:
		t.Fatalf("Timer should have been restarted upon second elapse")
	}

	timer.Elapse()
	select {
	case <-restartTimer:
		t.Fatalf("Timer should not have been restarted upon third elapse")
	default:
		// ok
	}
	select {
	case <-send:
		t.Fatalf("Controller re-sent when tries should have been exhausted")
	default:
		// ok
	}
}

func TestHandleAckRestartsTryLimit(t *testing.T) {
	resend := make(chan bool, 1)
	restartTimer := make(chan bool, 1)
	session := &readsession.MockReadSession{
		BeginHandler: func() {
		},
		ResendHandler: func() {
			resend <- true
		},
		HandleAckHandler: func(_ *safepackets.SafeAck) {
		},
	}
	timer := NewMockTimer(restartTimer, make(chan bool, 1))
	controller := manualTimeoutController(3, session, func() {}, timer)
	controller.BeginSession() // begin contains first try
	select {
	case <-restartTimer:
		// ok
	default:
		t.Fatalf("Timer should have been restarted upon begin")
	}

	timer.Elapse()
	select {
	case <-restartTimer:
		// ok
	default:
		t.Fatalf("Timer should have been restarted upon elapse")
	}
	select {
	case <-resend:
		// ok, second try
	default:
		t.Fatalf("Controller did not re-send in time")
	}

	timer.Elapse()
	runtime.Gosched()
	select {
	case <-restartTimer:
		// ok
	default:
		t.Fatalf("Timer should have been restarted upon elapse")
	}
	select {
	case <-resend:
		// ok, third try
	default:
		t.Fatalf("Controller did not re-send upon elapse")
	}

	controller.HandleAck(safepackets.NewSafeAck(5))
	for i := 0; i < 3; i++ {
		timer.Elapse()
		runtime.Gosched()
		select {
		case <-restartTimer:
			// ok
		default:
			t.Fatalf("Timer should have been restarted upon elapse")
		}
		select {
		case <-resend:
			// ok, try 1-2-3
		default:
			t.Fatalf("Controller did not re-send upon elapse")
		}
	}

	timer.Elapse()
	select {
	case <-resend:
		t.Fatalf("Controller re-sent when tries should have been exhausted")
	default:
		// ok, correct timeout
	}
}

func TestTimingOutWithOneTryCausesFinish(t *testing.T) {
	expired := make(chan bool, 1)
	session := &readsession.MockReadSession{
		BeginHandler: func() {
		},
	}
	timer := NewMockTimer(make(chan bool, 1), make(chan bool, 1))
	controller := manualTimeoutController(1, session, func() {
		expired <- true
	}, timer)
	controller.BeginSession()

	select {
	case <-expired:
	// ok
	default:
		t.Fatalf("Controller did not call expired callback after begin")
	}
}

func TestTimingOutWithMultipleTriesCausesFinish(t *testing.T) {
	begin := make(chan bool, 1)
	resend := make(chan bool, 1)
	restartTimer := make(chan bool, 1)
	expired := make(chan bool, 1)
	session := &readsession.MockReadSession{
		BeginHandler: func() {
			begin <- true
		},
		ResendHandler: func() {
			resend <- true
		},
	}
	timer := NewMockTimer(restartTimer, make(chan bool, 1))
	controller := manualTimeoutController(2, session, func() {
		expired <- true
	}, timer)
	controller.BeginSession()
	select {
	case <-restartTimer:
	// ok
	default:
		t.Fatalf("Timer not restarted after begin")
	}
	select {
	case <-begin:
	// ok
	default:
		t.Fatalf("Did not synchronize with begin")
	}

	timer.Elapse()
	select {
	case <-restartTimer:
	// ok
	default:
		t.Fatalf("Timer not restarted after elapse")
	}
	select {
	case <-resend:
	// ok
	case <-expired:
		t.Fatalf("Controller prematurely expired")
	default:
		t.Fatalf("Controller timed out too quickly")
	}

	timer.Elapse()
	runtime.Gosched()
	select {
	case <-restartTimer:
		t.Fatalf("Timer restarted after exhaustive elapse")
	default:
		// ok
	}

	select {
	case <-expired:
	// ok
	case <-resend:
		t.Fatalf("Controller resent when it should have expired")
	default:
		t.Fatalf("Controller did not call expired callback")
	}
}
