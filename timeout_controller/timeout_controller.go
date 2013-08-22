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
	c := &timeoutController{
		duration:       duration,
		tryLimit:       tryLimit,
		triesRemaining: tryLimit,
		timeout:        make(chan bool, tryLimit),
		timer:          time.NewTimer(duration),
	}

	c.timer.Stop()

	return c
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
	Stop()

	// Triggered when a countdown elapses.
	// Sends true when the countdown has expired (i.e. the maximum number of retries has been consumed).
	Timeout() chan bool
}

func (c *timeoutController) Countdown() error {
	if c.triesRemaining == 0 {
		return errors.New("No tries remaining")
	}

	go func() {
		c.timer.Reset(c.duration)
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

func (c *timeoutController) Stop() {
	c.timer.Stop()
}

func (c *timeoutController) Restart() {
	c.Stop()
	c.triesRemaining = c.tryLimit
}
