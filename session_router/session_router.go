package session_router

import (
	"github.com/mark-rushakoff/go_tftpd/readsessioncollection"
	"github.com/mark-rushakoff/go_tftpd/safetyfilter"
)

type SessionRouter struct {
	readSessions *readsessioncollection.ReadSessionCollection
}

func NewSessionRouter(readSessions *readsessioncollection.ReadSessionCollection) *SessionRouter {
	return &SessionRouter{
		readSessions: readSessions,
	}
}

func (r *SessionRouter) RouteAck(ack *safetyfilter.IncomingSafeAck) {
	session, found := r.readSessions.Fetch(ack.Addr)
	if !found {
		return
	}

	session.HandleAck(ack.Ack)
}
