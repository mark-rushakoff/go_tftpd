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
	config := &read_session.ReadSessionConfig{}
	m.readSessions[incomingReadRequest.Addr.String()] = m.readSessionFactory(config)
	m.lock.Unlock()
}
