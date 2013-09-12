package safety_filter

import (
	"net"
	"strings"

	"github.com/mark-rushakoff/go_tftpd/packets"
	"github.com/mark-rushakoff/go_tftpd/requestagent"
	"github.com/mark-rushakoff/go_tftpd/safepackets"
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
		Ack:  safepackets.NewSafeAck(incomingAck.Ack.BlockNumber),
	}
	f.Handler.HandleSafeAck(safeAck)
}

func (f *SafetyFilter) HandleReadRequest(incomingReadRequest *requestagent.IncomingReadRequest) {
	var mode safepackets.ReadWriteMode
	switch strings.ToLower(incomingReadRequest.Read.Mode) {
	case "netascii":
		mode = safepackets.NetAscii
		// TODO: handle and test octet case
	default:
		f.Handler.HandleError(&IncomingInvalidMessage{
			ErrorCode:    packets.Undefined,
			ErrorMessage: "Invalid mode string",
			Addr:         incomingReadRequest.Addr,
		})
		return
	}

	safeReadRequestPacket := &safepackets.SafeReadRequest{
		Filename: incomingReadRequest.Read.Filename,
		Mode:     mode,
	}

	safeReadRequest := &IncomingSafeReadRequest{
		Read: safeReadRequestPacket,
		Addr: incomingReadRequest.Addr,
	}
	f.Handler.HandleSafeReadRequest(safeReadRequest)
}
