package timeoutcontroller

import (
	"time"

	"github.com/mark-rushakoff/go_tftpd/readsession"
	"github.com/mark-rushakoff/go_tftpd/safepackets"
)

type TimeoutController interface {
	BeginSession()
	HandleAck(*safepackets.SafeAck)
}

type timeoutController struct {
	duration time.Duration

	tryCounter *tryCounter

	timer timer

	session readsession.ReadSession

	onExpire func()

	done chan bool
}

func NewTimeoutController(duration time.Duration, tryLimit uint, session readsession.ReadSession, onExpire func()) TimeoutController {
	timer := newTimer(duration)

	return manualTimeoutController(tryLimit, session, onExpire, timer)
}

func manualTimeoutController(tryLimit uint, session readsession.ReadSession, onExpire func(), timer timer) TimeoutController {
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

func (c *timeoutController) HandleAck(ack *safepackets.SafeAck) {
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
