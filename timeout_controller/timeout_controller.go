package timeout_controller

import (
	"time"

	"github.com/mark-rushakoff/go_tftpd/read_session"
	"github.com/mark-rushakoff/go_tftpd/safe_packets"
)

type TimeoutController struct {
	duration time.Duration
	tryLimit uint

	triesRemaining uint

	timer *time.Timer

	session read_session.ReadSession
}

func NewTimeoutController(duration time.Duration, tryLimit uint, session read_session.ReadSession) *TimeoutController {
	timer := time.NewTimer(time.Second)
	timer.Stop() // no false timeouts if there's a long time between initializing and calling Begin

	c := &TimeoutController{
		duration:       duration,
		tryLimit:       tryLimit,
		triesRemaining: tryLimit,
		timer:          timer,
		session:        session,
	}

	go func() {
		for {
			select {
			case <-c.timer.C:
				c.resendDueToTimeout()
				/* case <-done: */
				/* return */
			}
		}
	}()

	return c
}

func (c *TimeoutController) Begin() {
	c.session.Begin()
	c.triesRemaining--
	c.timer.Reset(c.duration)
}

func (c *TimeoutController) HandleAck(ack *safe_packets.SafeAck) {
	c.session.HandleAck(ack)
	c.triesRemaining = c.tryLimit
	c.timer.Reset(c.duration)
}

func (c *TimeoutController) resendDueToTimeout() {
	c.session.Resend()

	c.triesRemaining--
	if c.triesRemaining <= 0 {
	} else {
		c.timer.Reset(c.duration)
	}
}
