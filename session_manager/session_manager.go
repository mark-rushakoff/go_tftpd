package session_manager

import (
	"net"
	"sync"

	"github.com/mark-rushakoff/go_tftpd/read_session"
	"github.com/mark-rushakoff/go_tftpd/safety_filter"
)

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
	session := m.readSessionFactory(incomingReadRequest.Read.Filename)
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
