package session_manager

import (
	"sync"

	"github.com/mark-rushakoff/go_tftpd/read_session"
	"github.com/mark-rushakoff/go_tftpd/safety_filter"
)

type SessionManager struct {
	ReadRequest  chan *safety_filter.IncomingSafeReadRequest
	readSessions map[string]*read_session.ReadSession

	lock               sync.RWMutex
	readSessionFactory read_session.ReadSessionFactory
}

func NewSessionManager(readSessionFactory read_session.ReadSessionFactory) *SessionManager {
	return &SessionManager{
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
	m.readSessions[sessionKey(incomingReadRequest)] = session
	session.Begin()
}

func sessionKey(incomingReadRequest *safety_filter.IncomingSafeReadRequest) string {
	return incomingReadRequest.Addr.String()
}
