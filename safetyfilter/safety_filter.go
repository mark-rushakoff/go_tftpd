package safetyfilter

import (
	"net"

	"github.com/mark-rushakoff/go_tftpd/packets"
	"github.com/mark-rushakoff/go_tftpd/requestagent"
	"github.com/mark-rushakoff/go_tftpd/safepackets"
)

// SafetyFilter converts potentially unsafe messages from a RequestAgent into guaranteed-safe messages.
// Start it by calling Filter.
type SafetyFilter struct {
	converter safepackets.Converter
	handler   SafeRequestHandler
}

func MakeSafetyFilter(converter safepackets.Converter, handler SafeRequestHandler) *SafetyFilter {
	return &SafetyFilter{
		converter: converter,
		handler:   handler,
	}
}

type IncomingSafeAck struct {
	Ack  *safepackets.SafeAck
	Addr net.Addr
}

type IncomingSafeReadRequest struct {
	Read *safepackets.SafeReadRequest
	Addr net.Addr
}

type IncomingInvalidMessage struct {
	ErrorCode    packets.ErrorCode
	ErrorMessage string
	Addr         net.Addr
}

func (f *SafetyFilter) HandleAck(incomingAck *requestagent.IncomingAck) {
	safeAck := &IncomingSafeAck{
		Addr: incomingAck.Addr,
		Ack:  f.converter.FromAck(incomingAck.Ack),
	}
	f.handler.HandleSafeAck(safeAck)
}

func (f *SafetyFilter) HandleReadRequest(incomingReadRequest *requestagent.IncomingReadRequest) {
	safeReadRequestPacket, err := f.converter.FromReadRequest(incomingReadRequest.Read)
	if err != nil {
		f.handler.HandleError(&IncomingInvalidMessage{
			ErrorCode:    err.Code(),
			ErrorMessage: err.Error(),
			Addr:         incomingReadRequest.Addr,
		})
		return
	}

	safeReadRequest := &IncomingSafeReadRequest{
		Read: safeReadRequestPacket,
		Addr: incomingReadRequest.Addr,
	}
	f.handler.HandleSafeReadRequest(safeReadRequest)
}
