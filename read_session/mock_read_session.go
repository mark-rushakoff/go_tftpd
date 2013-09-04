package read_session

import (
	"github.com/mark-rushakoff/go_tftpd/safe_packets"
)

type MockReadSession struct {
	BeginHandler     func()
	HandleAckHandler func(ack *safe_packets.SafeAck)
	ResendHandler    func()
}

func (s *MockReadSession) Begin() {
	s.BeginHandler()
}

func (s *MockReadSession) HandleAck(ack *safe_packets.SafeAck) {
	s.HandleAckHandler(ack)
}

func (s *MockReadSession) Resend() {
	s.ResendHandler()
}
