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

func TestElapseDoublesUntilRestart(t *testing.T) {
	timer := newTimer(3 * time.Millisecond)
	timer.Restart()
	start := time.Now()

	var i uint
	for i = 0; i < 4; i++ {
		select {
		case <-timer.Elapsed():
			duration := time.Since(start)
			if duration < (3*(1<<i))*time.Millisecond {
				t.Fatalf("Elapsed too quickly (%v) on iteration %v", duration, i)
			} else if duration > (7*(1<<i))*time.Millisecond {
				t.Fatalf("Elapsed but took too long (%v) on iteration %v", duration, i)
			}
		case <-time.After((10 * (1 << i)) * time.Millisecond):
			t.Fatalf("Took too long to elapse on iteration %v", i)
		}
	}

	timer.Restart()
	start = time.Now()
	select {
	case <-timer.Elapsed():
		duration := time.Since(start)
		if duration < 3*time.Millisecond {
			t.Fatalf("Elapsed too quickly (%v) on iteration %v", duration, i)
		} else if duration > (7)*time.Millisecond {
			t.Fatalf("Elapsed but took too long (%v) on iteration %v", duration, i)
		}
	case <-time.After((10) * time.Millisecond):
		t.Fatalf("Took too long to elapse on iteration %v", i)
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
