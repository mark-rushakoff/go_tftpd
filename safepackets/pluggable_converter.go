package safepackets

import (
	"github.com/mark-rushakoff/go_tftpd/packets"
)

type PluggableConverter struct {
	FromAckHandler         func(ack *packets.Ack) *SafeAck
	FromReadRequestHandler func(read *packets.ReadRequest) (*SafeReadRequest, *ConversionError)
}

func (c *PluggableConverter) FromAck(ack *packets.Ack) *SafeAck {
	return c.FromAckHandler(ack)
}

func (c *PluggableConverter) FromReadRequest(read *packets.ReadRequest) (*SafeReadRequest, *ConversionError) {
	return c.FromReadRequestHandler(read)
}
