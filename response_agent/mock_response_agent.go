package response_agent

import (
	"sync"

	"github.com/mark-rushakoff/go_tftpd/safe_packets"
)

type MockResponseAgent struct {
	acks   []*safe_packets.SafeAck
	errors []*safe_packets.SafeError
	data   []*safe_packets.SafeData

	totalMessagesSent int

	lock sync.RWMutex
}

func MakeMockResponseAgent() *MockResponseAgent {
	return &MockResponseAgent{
		acks:   make([]*safe_packets.SafeAck, 5),
		errors: make([]*safe_packets.SafeError, 5),
		data:   make([]*safe_packets.SafeData, 5),
	}
}

func (a *MockResponseAgent) TotalMessagesSent() int {
	a.lock.RLock()
	defer a.lock.RUnlock()
	return a.totalMessagesSent
}

func (a *MockResponseAgent) MostRecentAck() *safe_packets.SafeAck {
	a.lock.Lock()
	defer a.lock.Unlock()
	if len(a.acks) == 0 {
		return nil
	}

	return a.acks[len(a.acks)-1]
}

func (a *MockResponseAgent) MostRecentData() *safe_packets.SafeData {
	a.lock.Lock()
	defer a.lock.Unlock()
	if len(a.data) == 0 {
		return nil
	}

	return a.data[len(a.data)-1]
}

func (a *MockResponseAgent) SendAck(ack *safe_packets.SafeAck) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.acks = append(a.acks, ack)
	a.totalMessagesSent++
}

func (a *MockResponseAgent) SendErrorPacket(e *safe_packets.SafeError) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.errors = append(a.errors, e)
	a.totalMessagesSent++
}

func (a *MockResponseAgent) SendData(data *safe_packets.SafeData) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.data = append(a.data, data)
	a.totalMessagesSent++
}

func (a *MockResponseAgent) Reset() {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.acks = a.acks[:0]
	a.errors = a.errors[:0]
	a.data = a.data[:0]
	a.totalMessagesSent = 0
}
