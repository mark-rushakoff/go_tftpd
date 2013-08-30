package safety_filter

type SafeRequestHandler interface {
	HandleSafeAck(*IncomingSafeAck)
	HandleSafeReadRequest(*IncomingSafeReadRequest)
}

type PluggableHandler struct {
	AckHandler         func(*IncomingSafeAck)
	ReadRequestHandler func(*IncomingSafeReadRequest)
}

func (h *PluggableHandler) HandleSafeAck(ack *IncomingSafeAck) {
	h.AckHandler(ack)
}

func (h *PluggableHandler) HandleSafeReadRequest(read *IncomingSafeReadRequest) {
	h.ReadRequestHandler(read)
}
