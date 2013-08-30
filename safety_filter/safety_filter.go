package safety_filter

import (
	"net"

	"github.com/mark-rushakoff/go_tftpd/request_agent"
	"github.com/mark-rushakoff/go_tftpd/safe_packets"
)

// SafetyFilter converts potentially unsafe messages from a RequestAgent into guaranteed-safe messages.
// Start it by calling Filter.
type SafetyFilter struct {
	Handler SafeRequestHandler
}

func MakeSafetyFilter(handler SafeRequestHandler) *SafetyFilter {
	return &SafetyFilter{
		Handler: handler,
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

func (f *SafetyFilter) HandleAck(incomingAck *request_agent.IncomingAck) {
	safeAck := &IncomingSafeAck{
		Addr: incomingAck.Addr,
		Ack:  safe_packets.NewSafeAck(incomingAck.Ack.BlockNumber),
	}
	f.Handler.HandleSafeAck(safeAck)
}

func (f *SafetyFilter) HandleReadRequest(incomingReadRequest *request_agent.IncomingReadRequest) {
	safeReadRequestPacket := &safe_packets.SafeReadRequest{
		Filename: incomingReadRequest.Read.Filename,
		Mode:     safe_packets.NetAscii, // TODO: handle octet and invalid cases
	}

	safeReadRequest := &IncomingSafeReadRequest{
		Read: safeReadRequestPacket,
		Addr: incomingReadRequest.Addr,
	}
	f.Handler.HandleSafeReadRequest(safeReadRequest)
}
