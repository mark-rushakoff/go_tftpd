package timeout_controller

import (
	"time"

	"github.com/mark-rushakoff/go_tftpd/read_session"
	"github.com/mark-rushakoff/go_tftpd/safe_packets"
)

type TimeoutController interface {
	Begin()
	HandleAck(*safe_packets.SafeAck)
}

type timeoutController struct {
	duration time.Duration
	tryLimit uint

	triesRemaining uint

	timer *time.Timer

	session read_session.ReadSession

	onExpire func()

	done chan bool
}

func NewTimeoutController(duration time.Duration, tryLimit uint, session read_session.ReadSession, onExpire func()) TimeoutController {
	timer := time.NewTimer(time.Second)
	timer.Stop() // no false timeouts if there's a long time between initializing and calling Begin

	c := &timeoutController{
		duration:       duration,
		tryLimit:       tryLimit,
		triesRemaining: tryLimit,
		timer:          timer,
		session:        session,
		onExpire:       onExpire,
		done:           make(chan bool),
	}

	go func() {
		for {
			select {
			case <-c.timer.C:
				c.resendDueToTimeout()
			case <-c.done:
				return
			}
		}
	}()

	return c
}

func (c *timeoutController) Begin() {
	c.session.Begin()
	c.triesRemaining--
	if c.triesRemaining > 0 {
		c.timer.Reset(c.duration)
	} else {
		c.expire()
	}
}

func (c *timeoutController) HandleAck(ack *safe_packets.SafeAck) {
	c.session.HandleAck(ack)
	c.triesRemaining = c.tryLimit
	c.timer.Reset(c.duration)
}

func (c *timeoutController) resendDueToTimeout() {
	if c.triesRemaining == 0 {
		c.expire()
		return
	}
	c.session.Resend()

	c.triesRemaining--
	c.timer.Reset(c.duration)
}

func (c *timeoutController) expire() {
	c.onExpire()
	c.done <- true
}
