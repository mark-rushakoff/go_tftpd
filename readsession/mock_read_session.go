package readsession

import (
	"github.com/mark-rushakoff/go_tftpd/safepackets"
)

type MockReadSession struct {
	BeginHandler     func()
	HandleAckHandler func(ack *safepackets.SafeAck)
	ResendHandler    func()
}

func (s *MockReadSession) Begin() {
	s.BeginHandler()
}

func (s *MockReadSession) HandleAck(ack *safepackets.SafeAck) {
	s.HandleAckHandler(ack)
}

func (s *MockReadSession) Resend() {
	s.ResendHandler()
}
