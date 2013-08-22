package read_session

import (
	"fmt"

	"github.com/mark-rushakoff/go_tftpd/safe_packets"
)

type ReadSession struct {
	config   ReadSessionConfig
	Ack      chan *safe_packets.SafeAck
	Finished chan bool

	currentBlockNumber uint16
	currentDataPacket  *safe_packets.SafeData
}

func NewReadSession(config ReadSessionConfig) *ReadSession {
	return &ReadSession{
		config:   config,
		Ack:      make(chan *safe_packets.SafeAck),
		Finished: make(chan bool),
	}
}

func (s *ReadSession) Begin() {
	s.config.ResponseAgent.SendAck(safe_packets.NewSafeAck(0))

	s.nextBlock()
	s.sendData()

	go s.watch()
}

func (s *ReadSession) watch() {
	for {
		select {
		case ack := <-s.Ack:
			if ack.BlockNumber == s.currentBlockNumber {
				s.nextBlock()
				s.sendData()
			} else if ack.BlockNumber == s.currentBlockNumber-1 {
				s.sendData()
			} else {
				panic(fmt.Sprintf("Could not handle received ack: %v", ack))
			}
		case isExpired := <-s.config.TimeoutController.Timeout():
			if isExpired {
				go func() {
					s.Finished <- true
				}()
				return
			} else {
				s.sendData()
			}
		}
	}
}

func (s *ReadSession) sendData() {
	s.config.ResponseAgent.SendData(s.currentDataPacket)
}

func (s *ReadSession) nextBlock() {
	s.currentBlockNumber++

	dataBytes := make([]byte, s.config.BlockSize)
	bytesRead, _ := s.config.Reader.Read(dataBytes) // TODO: be defensive
	dataBytes = dataBytes[:bytesRead]

	s.currentDataPacket = safe_packets.NewSafeData(s.currentBlockNumber, dataBytes)
}
