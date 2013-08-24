package request_agent

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"

	"github.com/mark-rushakoff/go_tftpd/packets"
)

type RequestAgent struct {
	Ack                 chan *IncomingAck
	Error               chan *IncomingError
	Data                chan *IncomingData
	ReadRequest         chan *IncomingReadRequest
	WriteRequest        chan *IncomingWriteRequest
	InvalidTransmission chan *InvalidTransmission
	conn                net.PacketConn
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

func NewRequestAgent(conn net.PacketConn) *RequestAgent {
	return &RequestAgent{
		Ack:                 make(chan *IncomingAck),
		Error:               make(chan *IncomingError),
		Data:                make(chan *IncomingData),
		ReadRequest:         make(chan *IncomingReadRequest),
		WriteRequest:        make(chan *IncomingWriteRequest),
		InvalidTransmission: make(chan *InvalidTransmission),
		conn:                conn,
	}
}

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
		go a.handleAck(b, addr)
	case packets.DataOpcode:
		go a.handleData(b, addr)
	case packets.ReadOpcode:
		go a.handleRead(b, addr)
	case packets.WriteOpcode:
		go a.handleWrite(b, addr)
	case packets.ErrorOpcode:
		go a.handleError(b, addr)
	default:
		go a.handleInvalidPacket(b, InvalidOpcode, addr)
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
	a.Ack <- &IncomingAck{&packets.Ack{blockNum}, addr}
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
	dataPacket := &packets.Data{blockNum, data}
	a.Data <- &IncomingData{dataPacket, addr}
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
	errorPacket := &packets.Error{packets.ErrorCode(code), message}
	a.Error <- &IncomingError{errorPacket, addr}
}

func (a *RequestAgent) handleRead(b []byte, addr net.Addr) {
	content, invalidReason, ok := parseReadWriteRequestContent(b)
	if ok {
		read := &packets.ReadRequest{content.filename, content.readWriteMode}
		a.ReadRequest <- &IncomingReadRequest{read, addr}
	} else {
		a.handleInvalidPacket(b, invalidReason, addr)
	}
}

func (a *RequestAgent) handleWrite(b []byte, addr net.Addr) {
	content, invalidReason, ok := parseReadWriteRequestContent(b)
	if ok {
		write := &packets.WriteRequest{content.filename, content.readWriteMode}
		a.WriteRequest <- &IncomingWriteRequest{write, addr}
	} else {
		a.handleInvalidPacket(b, invalidReason, addr)
	}
}

type readWriteRequestContent struct {
	filename      string
	readWriteMode string
}

func parseReadWriteRequestContent(b []byte) (content readWriteRequestContent, invalidReason InvalidTransmissionReason, ok bool) {
	remaining := b[2:]
	nulIndex := bytes.IndexByte(remaining, 0)
	if nulIndex == -1 {
		invalidReason = MissingField
		return
	}

	content.filename = string(remaining[:nulIndex])

	remaining = remaining[nulIndex+1:]
	nulIndex = bytes.IndexByte(remaining, 0)
	if nulIndex == -1 {
		invalidReason = MissingField
		return
	}

	if len(remaining) > nulIndex+1 {
		invalidReason = PacketTooLong
		return
	}

	content.readWriteMode = string(remaining[:nulIndex])

	ok = true
	return
}

func (a *RequestAgent) handleInvalidPacket(b []byte, reason InvalidTransmissionReason, addr net.Addr) {
	a.InvalidTransmission <- &InvalidTransmission{b, reason, addr}
}
