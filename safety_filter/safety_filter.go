package safety_filter

import (
	"net"

	"github.com/mark-rushakoff/go_tftpd/request_agent"
	"github.com/mark-rushakoff/go_tftpd/safe_packets"
)

// SafetyFilter converts potentially unsafe messages from a RequestAgent into guaranteed-safe messages.
// Start it by calling Filter.
type SafetyFilter struct {
	IncomingRead chan *IncomingSafeReadRequest
	IncomingAck  chan *IncomingSafeAck
	requestAgent *request_agent.RequestAgent
}

func MakeSafetyFilter(requestAgent *request_agent.RequestAgent) *SafetyFilter {
	filter := &SafetyFilter{
		IncomingAck:  make(chan *IncomingSafeAck),
		IncomingRead: make(chan *IncomingSafeReadRequest),
		requestAgent: requestAgent,
	}

	return filter
}

// Convert incoming requests from the associated requestAgent into safe requests that
// are output on the IncomingXXX channels.
func (f *SafetyFilter) Filter() {
	for {
		select {
		case incomingAck := <-f.requestAgent.Ack:
			safeAck := &IncomingSafeAck{
				Addr: incomingAck.Addr,
				Ack:  safe_packets.NewSafeAck(incomingAck.Ack.BlockNumber),
			}
			f.IncomingAck <- safeAck
		case incomingRead := <-f.requestAgent.ReadRequest:
			safeReadRequest := makeSafeReadRequest(incomingRead)
			f.IncomingRead <- safeReadRequest
		}
	}
}

type IncomingSafeAck struct {
	Ack  *safe_packets.SafeAck
	Addr net.Addr
}

type IncomingSafeReadRequest struct {
	Read *safe_packets.SafeReadRequest
	Addr net.Addr
}

func makeSafeReadRequest(incomingReadRequest *request_agent.IncomingReadRequest) *IncomingSafeReadRequest {
	safeReadRequest := &safe_packets.SafeReadRequest{
		Filename: incomingReadRequest.Read.Filename,
		Mode:     safe_packets.NetAscii, // TODO: handle octet and invalid cases
	}

	return &IncomingSafeReadRequest{
		Read: safeReadRequest,
		Addr: incomingReadRequest.Addr,
	}
}
