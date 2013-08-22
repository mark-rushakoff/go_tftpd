package timeout_controller

import (
	"testing"
	"time"
)

func TestCountdownTriggersTimeout(t *testing.T) {
	controller := NewTimeoutController(10*time.Millisecond, 2)
	err := controller.Countdown()
	if err != nil {
		t.Fatalf("Expected no errors, received %v", err)
	}
	select {
	case isExpired := <-controller.Timeout():
		if isExpired {
			t.Fatalf("Should not have expired yet")
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatalf("Did not receive signal in time")
	}

	err = controller.Countdown()
	if err != nil {
		t.Fatalf("Expected no errors, received %v", err)
	}
	select {
	case isExpired := <-controller.Timeout():
		if !isExpired {
			t.Fatalf("Should have expired now")
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatalf("Did not receive signal in time")
	}

	err = controller.Countdown()
	if err == nil {
		t.Fatalf("Expected an error when calling countdown when out of tries, but got nil")
	}
}

func TestStopCancelsCountdown(t *testing.T) {
	controller := NewTimeoutController(10*time.Millisecond, 2)
	err := controller.Countdown()
	if err != nil {
		t.Fatalf("Expected no errors, received %v", err)
	}

	time.Sleep(1 * time.Millisecond)
	controller.Stop()

	select {
	case <-controller.Timeout():
		t.Fatalf("Should not have received timeout")
	case <-time.After(20 * time.Millisecond):
		// success
	}
}

func TestStopBeforeCountdownReturnsError(t *testing.T) {
	controller := NewTimeoutController(10*time.Millisecond, 2)
	err := controller.Countdown()
	if err != nil {
		t.Fatalf("Expected no errors, received %v", err)
	}

	err = controller.Stop()
	if err == nil {
		t.Fatalf("Expected error but did not get one")
	}
}
