package requestagent

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"

	"github.com/mark-rushakoff/go_tftpd/packets"
)

// RequestAgent watches a PacketConn and emits potentially unsafe messages on its several exposed channels.
type RequestAgent struct {
	Handler RequestHandler
	conn    net.PacketConn
}

type IncomingAck struct {
	Ack  *packets.Ack
	Addr net.Addr
}

type IncomingError struct {
	Error *packets.Error
	Addr  net.Addr
}

type IncomingData struct {
	Data *packets.Data
	Addr net.Addr
}

type IncomingReadRequest struct {
	Read *packets.ReadRequest
	Addr net.Addr
}

type IncomingWriteRequest struct {
	Write *packets.WriteRequest
	Addr  net.Addr
}

func NewRequestAgent(conn net.PacketConn, handler RequestHandler) *RequestAgent {
	return &RequestAgent{
		conn: conn,

		Handler: handler,
	}
}

// Read a single message and emit it on the appropriate channel.
func (a *RequestAgent) Read() {
	const maxPacketSize = 516
	b := make([]byte, maxPacketSize)
	bytesRead, addr, err := a.conn.ReadFrom(b)
	if err != nil {
		panic(fmt.Sprintf("Error reading from connection: %v", err.Error()))
	}
	b = b[:bytesRead]

	if bytesRead < 3 {
		go a.handleInvalidPacket(b, PacketTooShort, addr)
		return
	}

	opcodeBuf := bytes.NewBuffer(b[0:2])
	var opcode uint16
	err = binary.Read(opcodeBuf, binary.BigEndian, &opcode)
	if err != nil {
		panic(fmt.Sprintf("Error while reading opcode from packet: %v", err.Error()))
	}

	switch opcode {
	case packets.AckOpcode:
		a.handleAck(b, addr)
	case packets.DataOpcode:
		a.handleData(b, addr)
	case packets.ReadOpcode:
		a.handleRead(b, addr)
	case packets.WriteOpcode:
		a.handleWrite(b, addr)
	case packets.ErrorOpcode:
		a.handleError(b, addr)
	default:
		a.handleInvalidPacket(b, InvalidOpcode, addr)
	}
}

func (a *RequestAgent) handleAck(b []byte, addr net.Addr) {
	if len(b) < 4 {
		a.handleInvalidPacket(b, PacketTooShort, addr)
		return
	}

	if len(b) > 4 {
		a.handleInvalidPacket(b, PacketTooLong, addr)
		return
	}

	blockNumBuf := bytes.NewBuffer(b[2:4])
	var blockNum uint16
	err := binary.Read(blockNumBuf, binary.BigEndian, &blockNum)
	if err != nil {
		panic(fmt.Sprintf("Error while reading blockNum from packet: %v", err))
	}
	a.Handler.HandleAck(&IncomingAck{
		Ack:  &packets.Ack{BlockNumber: blockNum},
		Addr: addr,
	})
}

func (a *RequestAgent) handleData(b []byte, addr net.Addr) {
	if len(b) < 4 {
		a.handleInvalidPacket(b, PacketTooShort, addr)
		return
	}

	blockNumBuf := bytes.NewBuffer(b[2:4])
	var blockNum uint16
	err := binary.Read(blockNumBuf, binary.BigEndian, &blockNum)
	if err != nil {
		panic(fmt.Sprintf("Error while reading blockNum from packet: %v", err))
	}
	data := b[4:]
	dataPacket := &packets.Data{BlockNumber: blockNum, Data: data}
	a.Handler.HandleData(&IncomingData{Data: dataPacket, Addr: addr})
}

func (a *RequestAgent) handleError(b []byte, addr net.Addr) {
	if len(b) < 4 {
		a.handleInvalidPacket(b, PacketTooShort, addr)
		return
	}

	codeBuf := bytes.NewBuffer(b[2:4])
	var code uint16
	err := binary.Read(codeBuf, binary.BigEndian, &code)
	if err != nil {
		panic(fmt.Sprintf("Error while reading code from packet: %v", err))
	}

	remaining := b[4:]
	nulIndex := bytes.IndexByte(remaining, 0)
	if nulIndex == -1 {
		a.handleInvalidPacket(b, MissingField, addr)
		return
	}

	if len(remaining) > nulIndex+1 {
		a.handleInvalidPacket(b, PacketTooLong, addr)
		return
	}

	message := string(remaining[:nulIndex])
	errorPacket := &packets.Error{Code: packets.ErrorCode(code), Message: message}
	a.Handler.HandleError(&IncomingError{errorPacket, addr})
}

func (a *RequestAgent) handleRead(b []byte, addr net.Addr) {
	content, invalidReason, ok := parseReadWriteRequestContent(b)
	if ok {
		read := &packets.ReadRequest{
			Filename: content.filename,
			Mode:     content.readWriteMode,
			Options:  content.options,
		}
		a.Handler.HandleReadRequest(&IncomingReadRequest{read, addr})
	} else {
		a.handleInvalidPacket(b, invalidReason, addr)
	}
}

func (a *RequestAgent) handleWrite(b []byte, addr net.Addr) {
	content, invalidReason, ok := parseReadWriteRequestContent(b)
	if ok {
		write := &packets.WriteRequest{
			Filename: content.filename,
			Mode:     content.readWriteMode,
			Options:  content.options,
		}
		a.Handler.HandleWriteRequest(&IncomingWriteRequest{write, addr})
	} else {
		a.handleInvalidPacket(b, invalidReason, addr)
	}
}

type readWriteRequestContent struct {
	filename      string
	readWriteMode string
	options       map[string]string
}

func parseReadWriteRequestContent(b []byte) (content readWriteRequestContent, invalidReason InvalidTransmissionReason, ok bool) {
	remaining := b[2:]
	nulIndex := bytes.IndexByte(remaining, 0)
	if nulIndex == -1 {
		invalidReason = MissingField
		return
	}

	content.filename = string(remaining[:nulIndex])
	content.options = make(map[string]string)

	remaining = remaining[nulIndex+1:]
	nulIndex = bytes.IndexByte(remaining, 0)
	if nulIndex == -1 {
		invalidReason = MissingField
		return
	}

	content.readWriteMode = string(remaining[:nulIndex])
	remaining = remaining[nulIndex+1:]

	if len(remaining) == 0 {
		ok = true
		return
	}

	for {
		keyNulIndex := bytes.IndexByte(remaining, 0)
		if keyNulIndex == -1 {
			invalidReason = OptionsMalformed
			return
		}
		key := string(remaining[:keyNulIndex])
		remaining = remaining[keyNulIndex+1:]

		valNulIndex := bytes.IndexByte(remaining, 0)
		if valNulIndex == -1 {
			invalidReason = OptionsMalformed
			return
		}
		val := string(remaining[:valNulIndex])
		remaining = remaining[valNulIndex+1:]

		content.options[key] = val

		if len(remaining) == 0 {
			ok = true
			return
		}
	}
}

func (a *RequestAgent) handleInvalidPacket(b []byte, reason InvalidTransmissionReason, addr net.Addr) {
	a.Handler.HandleInvalidTransmission(&InvalidTransmission{b, reason, addr})
}
