package request_router

import (
	"github.com/mark-rushakoff/go_tftpd/safety_filter"
	"github.com/mark-rushakoff/go_tftpd/session_manager"
)

// The RequestRouter handles the flow of data between a SafetyFilter and a SessionManager.
type RequestRouter struct {
	safetyFilter   *safety_filter.SafetyFilter
	sessionManager *session_manager.SessionManager
}

func NewRequestRouter(safetyFilter *safety_filter.SafetyFilter, sessionManager *session_manager.SessionManager) *RequestRouter {
	return &RequestRouter{
		safetyFilter:   safetyFilter,
		sessionManager: sessionManager,
	}
}

// Route reads from the safety filter forever and routes incoming data to the session manager.
func (r *RequestRouter) Route() {
	for {
		select {
		case incomingRead := <-r.safetyFilter.IncomingRead:
			r.sessionManager.ReadRequest <- incomingRead
		case incomingAck := <-r.safetyFilter.IncomingAck:
			r.sessionManager.Ack <- incomingAck
		}
	}
}
