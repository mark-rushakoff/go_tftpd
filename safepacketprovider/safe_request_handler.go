package safepacketprovider

import (
	"github.com/mark-rushakoff/go_tftpd/safety_filter"
)

type safeRequestHandler struct {
	safeAck            chan<- *safety_filter.IncomingSafeAck
	safeReadRequest    chan<- *safety_filter.IncomingSafeReadRequest
	safeInvalidMessage chan<- *safety_filter.IncomingInvalidMessage
}

func (h *safeRequestHandler) HandleSafeAck(a *safety_filter.IncomingSafeAck) {
	h.safeAck <- a
}

func (h *safeRequestHandler) HandleSafeReadRequest(r *safety_filter.IncomingSafeReadRequest) {
	h.safeReadRequest <- r
}

func (h *safeRequestHandler) HandleError(i *safety_filter.IncomingInvalidMessage) {
	h.safeInvalidMessage <- i
}
