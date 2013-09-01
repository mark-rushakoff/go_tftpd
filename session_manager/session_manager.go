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

func (m *SessionManager) MakeReadSessionFromIncomingRequest(incomingReadRequest *safety_filter.IncomingSafeReadRequest) {
	session := m.readSessionFactory(incomingReadRequest.Read.Filename, incomingReadRequest.Addr)
	key := sessionKey(incomingReadRequest.Addr)
	m.lock.Lock()
	m.readSessions[key] = session
	m.lock.Unlock()
	session.Begin()
}

func (m *SessionManager) SendAckToReadSession(incomingAck *safety_filter.IncomingSafeAck) {
	key := sessionKey(incomingAck.Addr)
	m.lock.RLock()
	session, found := m.readSessions[key]
	if !found {
		panic("Could not find a read session for address: " + key)
	}
	m.lock.RUnlock()
	session.Ack() <- incomingAck.Ack
}

func sessionKey(addr net.Addr) string {
	return addr.String()
}
