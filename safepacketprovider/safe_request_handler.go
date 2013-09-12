package safepacketprovider

import (
	"github.com/mark-rushakoff/go_tftpd/safetyfilter"
)

type safeRequestHandler struct {
	safeAck            chan<- *safetyfilter.IncomingSafeAck
	safeReadRequest    chan<- *safetyfilter.IncomingSafeReadRequest
	safeInvalidMessage chan<- *safetyfilter.IncomingInvalidMessage
}

func (h *safeRequestHandler) HandleSafeAck(a *safetyfilter.IncomingSafeAck) {
	h.safeAck <- a
}

func (h *safeRequestHandler) HandleSafeReadRequest(r *safetyfilter.IncomingSafeReadRequest) {
	h.safeReadRequest <- r
}

func (h *safeRequestHandler) HandleError(i *safetyfilter.IncomingInvalidMessage) {
	h.safeInvalidMessage <- i
}
