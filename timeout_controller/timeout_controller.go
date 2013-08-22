package timeout_controller

import (
	"errors"
	"time"
)

type timeoutController struct {
	duration time.Duration
	tryLimit uint

	triesRemaining uint

	timeout chan bool

	timer *time.Timer
}

func NewTimeoutController(duration time.Duration, tryLimit uint) TimeoutController {
	return &timeoutController{
		duration:       duration,
		tryLimit:       tryLimit,
		triesRemaining: tryLimit,
		timeout:        make(chan bool, tryLimit),
	}
}

type TimeoutController interface {
	// Starts a single session of the timeout.
	// Returns an error if no retries remain.
	Countdown() error

	// Begin an entirely new timeout cycle.
	// This resets both the current countdown timer and the number of retries consumed.
	Restart()

	// Disables any remaining countdown timer.
	// May only be followed by a call to Restart.
	Stop() error

	// Triggered when a countdown elapses.
	// Sends true when the countdown has expired (i.e. the maximum number of retries has been consumed).
	Timeout() chan bool
}

func (c *timeoutController) Countdown() error {
	if c.triesRemaining == 0 {
		return errors.New("No tries remaining")
	}

	go func() {
		c.timer = time.NewTimer(c.duration)
		select {
		case <-c.timer.C:
			c.triesRemaining--
			c.timeout <- (c.triesRemaining == 0)
		}
	}()

	return nil
}

func (c *timeoutController) Timeout() chan bool {
	return c.timeout
}

func (c *timeoutController) Stop() error {
	if c.timer == nil {
		return errors.New("tried to stop before countdown took effect")
	}

	c.timer.Stop()
	return nil
}

func (c *timeoutController) Restart() {
	c.Stop()
	c.triesRemaining = c.tryLimit
	c.timer = nil
}
