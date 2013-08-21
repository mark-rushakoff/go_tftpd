package response_agent

import (
	"github.com/mark-rushakoff/go_tftpd/safe_packets"
)

type MockResponseAgent struct {
	acks   []*safe_packets.SafeAck
	errors []*safe_packets.SafeError
	data   []*safe_packets.SafeData

	totalMessagesSent int
}

func MakeMockResponseAgent() *MockResponseAgent {
	return &MockResponseAgent{
		acks:   make([]*safe_packets.SafeAck, 5),
		errors: make([]*safe_packets.SafeError, 5),
		data:   make([]*safe_packets.SafeData, 5),
	}
}

func (a *MockResponseAgent) TotalMessagesSent() int {
	return a.totalMessagesSent
}

func (a *MockResponseAgent) MostRecentAck() *safe_packets.SafeAck {
	if len(a.acks) == 0 {
		return nil
	}

	return a.acks[len(a.acks)-1]
}

func (a *MockResponseAgent) MostRecentData() *safe_packets.SafeData {
	if len(a.data) == 0 {
		return nil
	}

	return a.data[len(a.data)-1]
}

func (a *MockResponseAgent) SendAck(ack *safe_packets.SafeAck) {
	a.acks = append(a.acks, ack)
	a.totalMessagesSent++
}

func (a *MockResponseAgent) SendErrorPacket(e *safe_packets.SafeError) {
	a.errors = append(a.errors, e)
	a.totalMessagesSent++
}

func (a *MockResponseAgent) SendData(data *safe_packets.SafeData) {
	a.data = append(a.data, data)
	a.totalMessagesSent++
}

func (a *MockResponseAgent) Reset() {
	a.acks = a.acks[:0]
	a.errors = a.errors[:0]
	a.data = a.data[:0]
	a.totalMessagesSent = 0
}
