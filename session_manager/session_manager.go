package session_manager

import (
	"net"
	"sync"

	"github.com/mark-rushakoff/go_tftpd/read_session"
	"github.com/mark-rushakoff/go_tftpd/safety_filter"
)

// Accepts incoming messages from a safety filter and creates read sessions or
// forward the requests to the appropriate read session.
type SessionManager struct {
	Ack          chan *safety_filter.IncomingSafeAck
	ReadRequest  chan *safety_filter.IncomingSafeReadRequest
	readSessions map[string]*read_session.ReadSession

	lock               sync.RWMutex
	readSessionFactory read_session.ReadSessionFactory
}

func NewSessionManager(readSessionFactory read_session.ReadSessionFactory) *SessionManager {
	return &SessionManager{
		Ack:          make(chan *safety_filter.IncomingSafeAck),
		ReadRequest:  make(chan *safety_filter.IncomingSafeReadRequest),
		readSessions: make(map[string]*read_session.ReadSession),

		readSessionFactory: readSessionFactory,
	}
}

// Block forever and handle requests from the associated SafetyFilter.
func (m *SessionManager) Watch() {
	for {
		select {
		case incomingReadRequest := <-m.ReadRequest:
			m.makeReadSessionFromIncomingRequest(incomingReadRequest)
		case incomingAck := <-m.Ack:
			m.sendAckToReadSession(incomingAck)
		}
	}
}

func (m *SessionManager) makeReadSessionFromIncomingRequest(incomingReadRequest *safety_filter.IncomingSafeReadRequest) {
	m.lock.Lock()
	defer m.lock.Unlock()
	session := m.readSessionFactory(incomingReadRequest.Read.Filename, incomingReadRequest.Addr)
	if session.Config == nil {
		panic("made session with a nil config")
	}
	m.readSessions[sessionKey(incomingReadRequest.Addr)] = session
	session.Begin()
}

func (m *SessionManager) sendAckToReadSession(incomingAck *safety_filter.IncomingSafeAck) {
	m.lock.Lock()
	defer m.lock.Unlock()
	session := m.readSessions[sessionKey(incomingAck.Addr)]
	go func() {
		session.Ack <- incomingAck.Ack
	}()
}

func sessionKey(addr net.Addr) string {
	return addr.String()
}
