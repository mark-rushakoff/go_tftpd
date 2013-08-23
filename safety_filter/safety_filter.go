package safety_filter

import (
	"net"

	"github.com/mark-rushakoff/go_tftpd/request_agent"
	"github.com/mark-rushakoff/go_tftpd/safe_packets"
)

type SafetyFilter struct {
	IncomingRead chan *IncomingSafeReadRequest
	IncomingAck  chan *IncomingSafeAck
}

func MakeSafetyFilter(requestAgent *request_agent.RequestAgent) *SafetyFilter {
	filter := &SafetyFilter{
		IncomingAck: make(chan *IncomingSafeAck),
	}

	go func() {
		for {
			select {
			case incomingAck := <-requestAgent.Ack:
				safeAck := &IncomingSafeAck{
					Addr: incomingAck.Addr,
					Ack:  safe_packets.NewSafeAck(incomingAck.Ack.BlockNumber),
				}
				filter.IncomingAck <- safeAck
			}
		}
	}()

	return filter
}

type IncomingSafeAck struct {
	Ack  *safe_packets.SafeAck
	Addr *net.Addr
}

type IncomingSafeReadRequest struct {
	Ack  *safe_packets.SafeReadRequest
	Addr *net.Addr
}
