package read_session

import (
	"io"
	"net"

	"github.com/mark-rushakoff/go_tftpd/safe_packets"
)

type ReadSessionFactory func(filename string, clientAddr net.Addr) *ReadSession

// Handles incoming messages from a single client and responds via the ResponseAgent inside a ReadSessionConfig.
type ReadSession struct {
	config   *ReadSessionConfig
	ack      chan *safe_packets.SafeAck
	finished chan bool

	currentBlockNumber uint16
	currentDataPacket  *safe_packets.SafeData
}

func NewReadSession(config *ReadSessionConfig) *ReadSession {
	return &ReadSession{
		config:   config,
		ack:      make(chan *safe_packets.SafeAck),
		finished: make(chan bool),
	}
}

func (s *ReadSession) Ack() chan<- *safe_packets.SafeAck {
	return s.ack
}

func (s *ReadSession) Finished() <-chan bool {
	return s.finished
}

func (s *ReadSession) Begin() {
	s.nextBlock()
	s.sendData()
	s.config.TimeoutController.Countdown()

	var isFinished bool
	for ; !isFinished; isFinished = s.watch() {
	}
}

func (s *ReadSession) watch() (isFinished bool) {
	select {
	case ack := <-s.ack:
		if ack.BlockNumber == s.currentBlockNumber {
			outOfBlocks := s.nextBlock()
			if outOfBlocks {
				s.finished <- true
				return true
			} else {
				s.sendData()
				s.config.TimeoutController.Restart()
				s.config.TimeoutController.Countdown()
			}
		} else if ack.BlockNumber == s.currentBlockNumber-1 {
			s.sendData()
			s.config.TimeoutController.Restart()
			s.config.TimeoutController.Countdown()
		} else {
			panic("A very old ack is currently undefined behavior")
		}
	case isExpired := <-s.config.TimeoutController.Timeout():
		if isExpired {
			s.finished <- true
			return true
		} else {
			s.sendData()
			s.config.TimeoutController.Countdown()
		}
	}

	return false
}

func (s *ReadSession) sendData() {
	s.config.ResponseAgent.SendData(s.currentDataPacket)
}

func (s *ReadSession) nextBlock() (isFinished bool) {
	dataBytes := make([]byte, s.config.BlockSize)
	if s.config.Reader == nil {
		panic("Config.Reader is nil")
	}

	bytesRead, err := s.config.Reader.Read(dataBytes)
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
