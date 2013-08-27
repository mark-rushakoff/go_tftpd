package read_session

import (
	"fmt"
	"io"

	"github.com/mark-rushakoff/go_tftpd/safe_packets"
)

type ReadSessionFactory func(string) *ReadSession

type ReadSession struct {
	Config   *ReadSessionConfig
	Ack      chan *safe_packets.SafeAck
	Finished chan bool

	currentBlockNumber uint16
	currentDataPacket  *safe_packets.SafeData
}

func NewReadSession(config ReadSessionConfig) *ReadSession {
	return &ReadSession{
		Config:   &config, // TODO: accept a reference instead
		Ack:      make(chan *safe_packets.SafeAck),
		Finished: make(chan bool),
	}
}

func (s *ReadSession) Begin() {
	if s.Config == nil {
		panic("nil config")
	}

	if s.Config.ResponseAgent == nil {
		panic("nil response agent")
	}

	s.Config.ResponseAgent.SendAck(safe_packets.NewSafeAck(0))

	s.nextBlock()
	s.sendData()

	go s.watch()
}

func (s *ReadSession) watch() {
	if s.Config.TimeoutController == nil {
		panic("TimeoutController is nil")
	}

	for {
		select {
		case ack := <-s.Ack:
			if ack.BlockNumber == s.currentBlockNumber {
				isFinished := s.nextBlock()
				if isFinished {
					s.Finished <- true
				} else {
					s.sendData()
				}
			} else if ack.BlockNumber == s.currentBlockNumber-1 {
				s.sendData()
			} else {
				panic(fmt.Sprintf("Could not handle received ack: %v", ack))
			}
		case isExpired := <-s.Config.TimeoutController.Timeout():
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
	s.Config.ResponseAgent.SendData(s.currentDataPacket)
}

func (s *ReadSession) nextBlock() (isFinished bool) {
	dataBytes := make([]byte, s.Config.BlockSize)
	if s.Config.Reader == nil {
		panic("Config.Reader is nil")
	}

	bytesRead, err := s.Config.Reader.Read(dataBytes)
	if bytesRead == 0 {
		if err == io.EOF {
			return true
		} else {
			panic("Not sure what to do with a non-eof io error and 0 bytes read")
		}
	}

	dataBytes = dataBytes[:bytesRead]

	s.currentBlockNumber++
	s.currentDataPacket = safe_packets.NewSafeData(s.currentBlockNumber, dataBytes)
	return false
}
