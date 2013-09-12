package timeoutcontroller

import (
	"testing"
	"time"
)

func TestRestartSendsToElapsedWhenFinished(t *testing.T) {
	timer := newTimer(3 * time.Millisecond)

	timer.Restart()

	time.Sleep(2 * time.Millisecond)
	select {
	case <-timer.Elapsed():
		t.Fatalf("Timer elapsed too early")
	default:
		// ok
	}

	select {
	case <-timer.Elapsed():
		// ok
	case <-time.After(2 * time.Millisecond):
		t.Fatalf("Timer did not elapse as expected")
	}
}

func TestRestartRestartsTimer(t *testing.T) {
	timer := newTimer(10 * time.Millisecond)

	timer.Restart()

	time.Sleep(7 * time.Millisecond)
	select {
	case <-timer.Elapsed():
		t.Fatalf("Timer elapsed too early")
	default:
		// ok
	}

	timer.Restart()
	select {
	case <-timer.Elapsed():
		t.Fatalf("Timer should not have elapsed yet")
	case <-time.After(6 * time.Millisecond):
		// ok
	}

	select {
	case <-timer.Elapsed():
		// ok
	case <-time.After(5 * time.Millisecond):
		t.Fatalf("Timer should have elapsed by now")
	}
}

func TestDestroyDoesNotElapseTimer(t *testing.T) {
	timer := newTimer(3 * time.Millisecond)

	timer.Restart()

	time.Sleep(2 * time.Millisecond)
	select {
	case <-timer.Elapsed():
		t.Fatalf("Timer elapsed too early")
	default:
		// ok
	}

	timer.Destroy()
	select {
	case <-timer.Elapsed():
		t.Fatalf("Timer should not have elapsed after destroy")
	case <-time.After(4 * time.Millisecond):
		// ok
	}
}
