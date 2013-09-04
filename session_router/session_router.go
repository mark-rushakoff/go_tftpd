package session_router

import (
	"github.com/mark-rushakoff/go_tftpd/read_session_collection"
	"github.com/mark-rushakoff/go_tftpd/safety_filter"
)

type SessionRouter struct {
	readSessions *read_session_collection.ReadSessionCollection
}

func NewSessionRouter(readSessions *read_session_collection.ReadSessionCollection) *SessionRouter {
	return &SessionRouter{
		readSessions: readSessions,
	}
}

func (r *SessionRouter) RouteAck(ack *safety_filter.IncomingSafeAck) {
	session, found := r.readSessions.Fetch(ack.Addr)
	if !found {
		panic("Tried to route ack to unknown session")
	}

	session.HandleAck(ack.Ack)
}
