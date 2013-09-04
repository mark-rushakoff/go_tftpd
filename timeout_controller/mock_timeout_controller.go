package timeout_controller

import (
	"github.com/mark-rushakoff/go_tftpd/safe_packets"
)

type MockTimeoutController struct {
	HandleAckHandler func(*safe_packets.SafeAck)
	BeginHandler     func()
}

func (c *MockTimeoutController) HandleAck(ack *safe_packets.SafeAck) {
	c.HandleAckHandler(ack)
}

func (c *MockTimeoutController) Begin() {
	c.BeginHandler()
}
