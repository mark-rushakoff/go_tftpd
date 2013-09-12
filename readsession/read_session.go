package readsession

import (
	"io"

	"github.com/mark-rushakoff/go_tftpd/safe_packets"
)

type Config struct {
	Reader    io.Reader
	BlockSize uint16
}

type ReadSession interface {
	Begin()
	HandleAck(ack *safe_packets.SafeAck)
	Resend()
}

type readSession struct {
	config  *Config
	handler OutgoingHandler

	currentBlockNumber uint16
	currentDataPacket  *safe_packets.SafeData

	dataExhausted bool
	onFinish      func()
}

func NewReadSession(config *Config, handler OutgoingHandler, onFinish func()) *readSession {
	return &readSession{
		config:   config,
		handler:  handler,
		onFinish: onFinish,
	}
}

func (s *readSession) Begin() {
	s.nextBlock()
	s.sendData()
}

func (s *readSession) HandleAck(ack *safe_packets.SafeAck) {
	if s.dataExhausted {
		s.onFinish()
		return
	}

	if ack.BlockNumber == s.currentBlockNumber {
		s.nextBlock()
		s.sendData()
	} else if ack.BlockNumber == s.currentBlockNumber-1 {
		s.sendData()
	} else {
		panic("A very old ack is currently undefined behavior")
	}
}

func (s *readSession) Resend() {
	s.sendData()
}

func (s *readSession) sendData() {
	s.handler.SendData(s.currentDataPacket)
}

func (s *readSession) nextBlock() {
	dataBytes := make([]byte, s.config.BlockSize)
	if s.config.Reader == nil {
		panic("Config.Reader is nil")
	}

	bytesRead, err := s.config.Reader.Read(dataBytes)
	if bytesRead == 0 {
		s.dataExhausted = true
		if err == io.EOF {
			return
		} else {
			panic("Not sure what to do with a non-eof io error and 0 bytes read")
		}
	}

	if bytesRead < len(dataBytes) {
		s.dataExhausted = true
	}

	dataBytes = dataBytes[:bytesRead]

	s.currentBlockNumber++
	s.currentDataPacket = safe_packets.NewSafeData(s.currentBlockNumber, dataBytes)
}
