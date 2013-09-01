package safe_packet_provider

import (
	"github.com/mark-rushakoff/go_tftpd/request_agent"
	"github.com/mark-rushakoff/go_tftpd/safety_filter"
)

type requestHandler struct {
	safetyFilter *safety_filter.SafetyFilter
}

func (h *requestHandler) HandleAck(a *request_agent.IncomingAck) {
	h.safetyFilter.HandleAck(a)
}
func (h *requestHandler) HandleError(e *request_agent.IncomingError) {
	panic("idk how to handle error")
}
func (h *requestHandler) HandleData(d *request_agent.IncomingData) {
	panic("idk how to handle data")
}
func (h *requestHandler) HandleReadRequest(r *request_agent.IncomingReadRequest) {
	h.safetyFilter.HandleReadRequest(r)
}
func (h *requestHandler) HandleWriteRequest(w *request_agent.IncomingWriteRequest) {
	panic("idk how to handle write")
}
func (h *requestHandler) HandleInvalidTransmission(t *request_agent.InvalidTransmission) {
	panic("idk how to handle invalid tx")
}
