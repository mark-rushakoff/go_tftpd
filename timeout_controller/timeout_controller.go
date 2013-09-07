package timeout_controller

import (
	"time"

	"github.com/mark-rushakoff/go_tftpd/read_session"
	"github.com/mark-rushakoff/go_tftpd/safe_packets"
)

type TimeoutController interface {
	BeginSession()
	HandleAck(*safe_packets.SafeAck)
}

type timeoutController struct {
	duration time.Duration

	tryCounter *tryCounter

	timer timer

	session read_session.ReadSession

	onExpire func()

	done chan bool
}

func NewTimeoutController(duration time.Duration, tryLimit uint, session read_session.ReadSession, onExpire func()) TimeoutController {
	timer := newTimer(duration)

	return manualTimeoutController(tryLimit, session, onExpire, timer)
}

func manualTimeoutController(tryLimit uint, session read_session.ReadSession, onExpire func(), timer timer) TimeoutController {
	counter := &tryCounter{
		tryLimit:       tryLimit,
		triesRemaining: tryLimit,
	}

	c := &timeoutController{
		tryCounter: counter,
		timer:      timer,
		session:    session,
		onExpire:   onExpire,
		done:       make(chan bool),
	}

	go func() {
		for {
			select {
			case <-timer.Elapsed():
				c.resendDueToTimeout()
			case <-c.done:
				return
			}
		}
	}()

	return c
}

func (c *timeoutController) BeginSession() {
	c.session.Begin()
	c.tryCounter.Decrement()
	if c.tryCounter.IsZero() {
		c.expire()
	} else {
		c.timer.Restart()
	}
}

func (c *timeoutController) HandleAck(ack *safe_packets.SafeAck) {
	c.session.HandleAck(ack)
	c.tryCounter.Reset()
	c.timer.Restart()
}

func (c *timeoutController) resendDueToTimeout() {
	if c.tryCounter.IsZero() {
		c.expire()
		return
	}
	c.session.Resend()

	c.tryCounter.Decrement()
	c.timer.Restart()
}

func (c *timeoutController) expire() {
	c.onExpire()
	c.done <- true
	c.timer.Destroy()
}
