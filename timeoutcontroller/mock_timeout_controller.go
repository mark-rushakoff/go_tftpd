package timeoutcontroller

import (
	"github.com/mark-rushakoff/go_tftpd/safepackets"
)

type MockTimeoutController struct {
	HandleAckHandler    func(*safepackets.SafeAck)
	BeginSessionHandler func()
}

func (c *MockTimeoutController) HandleAck(ack *safepackets.SafeAck) {
	c.HandleAckHandler(ack)
}

func (c *MockTimeoutController) BeginSession() {
	c.BeginSessionHandler()
}
