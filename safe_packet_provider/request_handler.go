package safe_packet_provider

import (
	"github.com/mark-rushakoff/go_tftpd/requestagent"
	"github.com/mark-rushakoff/go_tftpd/safety_filter"
)

type requestHandler struct {
	safetyFilter *safety_filter.SafetyFilter
}

func (h *requestHandler) HandleAck(a *requestagent.IncomingAck) {
	h.safetyFilter.HandleAck(a)
}
func (h *requestHandler) HandleError(e *requestagent.IncomingError) {
	panic("idk how to handle error")
}
func (h *requestHandler) HandleData(d *requestagent.IncomingData) {
	panic("idk how to handle data")
}
func (h *requestHandler) HandleReadRequest(r *requestagent.IncomingReadRequest) {
	h.safetyFilter.HandleReadRequest(r)
}
func (h *requestHandler) HandleWriteRequest(w *requestagent.IncomingWriteRequest) {
	panic("idk how to handle write")
}
func (h *requestHandler) HandleInvalidTransmission(t *requestagent.InvalidTransmission) {
	panic("idk how to handle invalid tx")
}
