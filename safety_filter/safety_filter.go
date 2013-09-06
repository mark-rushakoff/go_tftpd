package safety_filter

import (
	"net"
	"strings"

	"github.com/mark-rushakoff/go_tftpd/packets"
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

type IncomingInvalidMessage struct {
	ErrorCode    packets.ErrorCode
	ErrorMessage string
	Addr         net.Addr
}

func (f *SafetyFilter) HandleAck(incomingAck *request_agent.IncomingAck) {
	safeAck := &IncomingSafeAck{
		Addr: incomingAck.Addr,
		Ack:  safe_packets.NewSafeAck(incomingAck.Ack.BlockNumber),
	}
	f.Handler.HandleSafeAck(safeAck)
}

func (f *SafetyFilter) HandleReadRequest(incomingReadRequest *request_agent.IncomingReadRequest) {
	var mode safe_packets.ReadWriteMode
	switch strings.ToLower(incomingReadRequest.Read.Mode) {
	case "netascii":
		mode = safe_packets.NetAscii
		// TODO: handle and test octet case
	default:
		f.Handler.HandleError(&IncomingInvalidMessage{
			ErrorCode:    packets.Undefined,
			ErrorMessage: "Invalid mode string",
			Addr:         incomingReadRequest.Addr,
		})
		return
	}

	safeReadRequestPacket := &safe_packets.SafeReadRequest{
		Filename: incomingReadRequest.Read.Filename,
		Mode:     mode,
	}

	safeReadRequest := &IncomingSafeReadRequest{
		Read: safeReadRequestPacket,
		Addr: incomingReadRequest.Addr,
	}
	f.Handler.HandleSafeReadRequest(safeReadRequest)
}
