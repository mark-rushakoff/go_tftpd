package timeout_controller

import (
	"testing"
	"time"
)

func TestResetSendsToElapsedWhenFinished(t *testing.T) {
	timer := newTimer(3 * time.Millisecond)

	timer.Reset()

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

func TestResetRestartsTimer(t *testing.T) {
	timer := newTimer(3 * time.Millisecond)

	timer.Reset()

	time.Sleep(2 * time.Millisecond)
	select {
	case <-timer.Elapsed():
		t.Fatalf("Timer elapsed too early")
	default:
		// ok
	}

	timer.Reset()
	select {
	case <-timer.Elapsed():
		t.Fatalf("Timer should not have elapsed yet")
	case <-time.After(2 * time.Millisecond):
		// ok
	}

	select {
	case <-timer.Elapsed():
		// ok
	case <-time.After(2 * time.Millisecond):
		t.Fatalf("Timer should have elapsed by now")
	}
}

func TestDestroyDoesNotElapseTimer(t *testing.T) {
	timer := newTimer(3 * time.Millisecond)

	timer.Reset()

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
