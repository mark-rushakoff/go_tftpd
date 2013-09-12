package readsessioncollection

import (
	"net"
	"sync"

	"github.com/mark-rushakoff/go_tftpd/timeoutcontroller"
)

type sessionKey string

type ReadSessionCollection struct {
	sessions map[sessionKey]timeoutcontroller.TimeoutController
	lock     sync.RWMutex
}

func NewReadSessionCollection() *ReadSessionCollection {
	return &ReadSessionCollection{
		sessions: make(map[sessionKey]timeoutcontroller.TimeoutController),
	}
}

func (s *ReadSessionCollection) Add(session timeoutcontroller.TimeoutController, addr net.Addr) {
	key := key(addr)
	s.add(session, key)
}

func (s *ReadSessionCollection) Fetch(addr net.Addr) (session timeoutcontroller.TimeoutController, ok bool) {
	return s.fetch(key(addr))
}

func (s *ReadSessionCollection) Remove(addr net.Addr) {
	s.remove(key(addr))
}

func (s *ReadSessionCollection) add(session timeoutcontroller.TimeoutController, key sessionKey) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.sessions[key] = session
}

func (s *ReadSessionCollection) fetch(key sessionKey) (session timeoutcontroller.TimeoutController, ok bool) {
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
