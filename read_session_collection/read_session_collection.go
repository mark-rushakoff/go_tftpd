package read_session_collection

import (
	"net"
	"sync"

	"github.com/mark-rushakoff/go_tftpd/timeout_controller"
)

type sessionKey string

type ReadSessionCollection struct {
	sessions map[sessionKey]timeout_controller.TimeoutController
	lock     sync.RWMutex
}

func NewReadSessionCollection() *ReadSessionCollection {
	return &ReadSessionCollection{
		sessions: make(map[sessionKey]timeout_controller.TimeoutController),
	}
}

func (s *ReadSessionCollection) Add(session timeout_controller.TimeoutController, addr net.Addr) {
	key := key(addr)
	s.add(session, key)
}

func (s *ReadSessionCollection) Fetch(addr net.Addr) (session timeout_controller.TimeoutController, ok bool) {
	return s.fetch(key(addr))
}

func (s *ReadSessionCollection) Remove(addr net.Addr) {
	s.remove(key(addr))
}

func (s *ReadSessionCollection) add(session timeout_controller.TimeoutController, key sessionKey) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.sessions[key] = session
}

func (s *ReadSessionCollection) fetch(key sessionKey) (session timeout_controller.TimeoutController, ok bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	session, ok = s.sessions[key]
	return
}

func (s *ReadSessionCollection) remove(key sessionKey) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.sessions, key)
}

func key(addr net.Addr) sessionKey {
	return sessionKey(addr.String())
}
