package read_session

import (
	"github.com/mark-rushakoff/go_tftpd/safe_packets"
)

type ReadSession struct {
	config ReadSessionConfig
}

func (s *ReadSession) Begin() {
	s.config.ResponseAgent.SendAck(safe_packets.NewSafeAck(0))

	dataBytes := make([]byte, s.config.BlockSize)
	bytesRead, _ := s.config.Reader.Read(dataBytes) // TODO: be defensive
	dataBytes = dataBytes[:bytesRead]
	s.config.ResponseAgent.SendData(safe_packets.NewSafeData(1, dataBytes))
}
