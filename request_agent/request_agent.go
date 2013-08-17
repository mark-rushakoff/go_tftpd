package request_agent

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strings"

	"github.com/mark-rushakoff/go_tftpd/messages"
)

type RequestAgent struct {
	conn         net.PacketConn
	Ack          chan *messages.Ack
	Error        chan *messages.Error
	Data         chan *messages.Data
	ReadRequest  chan *messages.ReadRequest
	WriteRequest chan *messages.WriteRequest
}

func NewRequestAgent(conn net.PacketConn) *RequestAgent {
	return &RequestAgent{
		conn:         conn,
		Ack:          make(chan *messages.Ack),
		Error:        make(chan *messages.Error),
		Data:         make(chan *messages.Data),
		ReadRequest:  make(chan *messages.ReadRequest),
		WriteRequest: make(chan *messages.WriteRequest),
	}
}

func (a *RequestAgent) Read() {
	const maxPacketSize = 516
	b := make([]byte, maxPacketSize)
	bytesRead, _, _ := a.conn.ReadFrom(b)
	b = b[:bytesRead]
	opcodeBuf := bytes.NewBuffer(b[0:2])
	var opcode uint16
	err := binary.Read(opcodeBuf, binary.BigEndian, &opcode)
	if err != nil {
		panic("Error while reading opcode from packet")
	}

	switch opcode {
	case messages.AckOpcode:
		go a.handleAck(b)
	case messages.DataOpcode:
		go a.handleData(b)
	case messages.ReadOpcode:
		go a.handleRead(b)
	case messages.ErrorOpcode:
		go a.handleError(b)
	default:
		panic(fmt.Sprintf("Unknown opcode: %v", opcode))
	}
}

func (a *RequestAgent) handleAck(b []byte) {
	blockNumBuf := bytes.NewBuffer(b[2:4])
	var blockNum uint16
	err := binary.Read(blockNumBuf, binary.BigEndian, &blockNum)
	if err != nil {
		panic(fmt.Sprintf("Error while reading blockNum from packet: %v", err))
	}
	a.Ack <- &messages.Ack{blockNum}
}

func (a *RequestAgent) handleData(b []byte) {
	blockNumBuf := bytes.NewBuffer(b[2:4])
	var blockNum uint16
	err := binary.Read(blockNumBuf, binary.BigEndian, &blockNum)
	if err != nil {
		panic(fmt.Sprintf("Error while reading blockNum from packet: %v", err))
	}
	data := b[4:]
	a.Data <- &messages.Data{blockNum, data}
}

func (a *RequestAgent) handleError(b []byte) {
	codeBuf := bytes.NewBuffer(b[2:4])
	var code uint16
	err := binary.Read(codeBuf, binary.BigEndian, &code)
	if err != nil {
		panic(fmt.Sprintf("Error while reading code from packet: %v", err))
	}

	message := string(b[4 : len(b)-1]) // chop nul byte at end
	a.Error <- &messages.Error{messages.ErrorCode(code), message}
}

func (a *RequestAgent) handleRead(b []byte) {
	remaining := b[2:]
	nulIndex := bytes.IndexByte(remaining, 0)
	if nulIndex == -1 {
		panic("Could not find nul-terminator in Read request")
	}
	fileName := string(remaining[:nulIndex])

	modeStr := string(remaining[nulIndex+1 : len(remaining)-1])
	if strings.ToLower(modeStr) == "netascii" {
		a.ReadRequest <- &messages.ReadRequest{fileName, messages.NetAscii}
	} else {
		panic(fmt.Sprintf("Unrecognized mode in read request: %v", modeStr))
	}
}
