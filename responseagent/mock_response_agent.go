package responseagent

import (
	"sync"

	"github.com/mark-rushakoff/go_tftpd/safepackets"
)

type MockResponseAgent struct {
	acks   []*safepackets.SafeAck
	errors []*safepackets.SafeError
	data   []*safepackets.SafeData

	totalMessagesSent int

	lock sync.RWMutex
}

func MakeMockResponseAgent() *MockResponseAgent {
	return &MockResponseAgent{
		acks:   make([]*safepackets.SafeAck, 5),
		errors: make([]*safepackets.SafeError, 5),
		data:   make([]*safepackets.SafeData, 5),
	}
}

func (a *MockResponseAgent) TotalMessagesSent() int {
	a.lock.RLock()
	defer a.lock.RUnlock()
	return a.totalMessagesSent
}

func (a *MockResponseAgent) MostRecentAck() *safepackets.SafeAck {
	a.lock.Lock()
	defer a.lock.Unlock()
	if len(a.acks) == 0 {
		return nil
	}

	return a.acks[len(a.acks)-1]
}

func (a *MockResponseAgent) MostRecentData() *safepackets.SafeData {
	a.lock.Lock()
	defer a.lock.Unlock()
	if len(a.data) == 0 {
		return nil
	}

	return a.data[len(a.data)-1]
}

func (a *MockResponseAgent) SendAck(ack *safepackets.SafeAck) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.acks = append(a.acks, ack)
	a.totalMessagesSent++
}

func (a *MockResponseAgent) SendErrorPacket(e *safepackets.SafeError) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.errors = append(a.errors, e)
	a.totalMessagesSent++
}

func (a *MockResponseAgent) SendData(data *safepackets.SafeData) {
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
