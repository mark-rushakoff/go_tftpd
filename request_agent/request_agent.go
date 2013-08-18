package request_agent

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"

	"github.com/mark-rushakoff/go_tftpd/packets"
)

type RequestAgent struct {
	Ack                 chan *packets.Ack
	Error               chan *packets.Error
	Data                chan *packets.Data
	ReadRequest         chan *packets.ReadRequest
	WriteRequest        chan *packets.WriteRequest
	InvalidTransmission chan *InvalidTransmission
	conn                net.PacketConn
}

func NewRequestAgent(conn net.PacketConn) *RequestAgent {
	return &RequestAgent{
		Ack:                 make(chan *packets.Ack),
		Error:               make(chan *packets.Error),
		Data:                make(chan *packets.Data),
		ReadRequest:         make(chan *packets.ReadRequest),
		WriteRequest:        make(chan *packets.WriteRequest),
		InvalidTransmission: make(chan *InvalidTransmission),
		conn:                conn,
	}
}

func (a *RequestAgent) Read() {
	const maxPacketSize = 516
	b := make([]byte, maxPacketSize)
	bytesRead, _, _ := a.conn.ReadFrom(b)
	b = b[:bytesRead]

	if bytesRead < 3 {
		go a.handleInvalidPacket(b, PacketTooShort)
		return
	}

	opcodeBuf := bytes.NewBuffer(b[0:2])
	var opcode uint16
	err := binary.Read(opcodeBuf, binary.BigEndian, &opcode)
	if err != nil {
		panic("Error while reading opcode from packet")
	}

	switch opcode {
	case packets.AckOpcode:
		go a.handleAck(b)
	case packets.DataOpcode:
		go a.handleData(b)
	case packets.ReadOpcode:
		go a.handleRead(b)
	case packets.WriteOpcode:
		go a.handleWrite(b)
	case packets.ErrorOpcode:
		go a.handleError(b)
	default:
		go a.handleInvalidPacket(b, InvalidOpcode)
	}
}

func (a *RequestAgent) handleAck(b []byte) {
	blockNumBuf := bytes.NewBuffer(b[2:4])
	var blockNum uint16
	err := binary.Read(blockNumBuf, binary.BigEndian, &blockNum)
	if err != nil {
		panic(fmt.Sprintf("Error while reading blockNum from packet: %v", err))
	}
	a.Ack <- &packets.Ack{blockNum}
}

func (a *RequestAgent) handleData(b []byte) {
	blockNumBuf := bytes.NewBuffer(b[2:4])
	var blockNum uint16
	err := binary.Read(blockNumBuf, binary.BigEndian, &blockNum)
	if err != nil {
		panic(fmt.Sprintf("Error while reading blockNum from packet: %v", err))
	}
	data := b[4:]
	a.Data <- &packets.Data{blockNum, data}
}

func (a *RequestAgent) handleError(b []byte) {
	codeBuf := bytes.NewBuffer(b[2:4])
	var code uint16
	err := binary.Read(codeBuf, binary.BigEndian, &code)
	if err != nil {
		panic(fmt.Sprintf("Error while reading code from packet: %v", err))
	}

	message := string(b[4 : len(b)-1]) // chop nul byte at end
	a.Error <- &packets.Error{packets.ErrorCode(code), message}
}

func (a *RequestAgent) handleRead(b []byte) {
	content, invalidReason, ok := parseReadWriteRequestContent(b)
	if ok {
		a.ReadRequest <- &packets.ReadRequest{content.filename, content.readWriteMode}
	} else {
		a.handleInvalidPacket(b, invalidReason)
	}
}

func (a *RequestAgent) handleWrite(b []byte) {
	content, invalidReason, ok := parseReadWriteRequestContent(b)
	if ok {
		a.WriteRequest <- &packets.WriteRequest{content.filename, content.readWriteMode}
	} else {
		a.handleInvalidPacket(b, invalidReason)
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

	content.readWriteMode = string(remaining[:nulIndex])

	ok = true
	return
}

func (a *RequestAgent) handleInvalidPacket(b []byte, reason InvalidTransmissionReason) {
	a.InvalidTransmission <- &InvalidTransmission{b, reason}
}
