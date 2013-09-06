package safety_filter

type SafeRequestHandler interface {
	HandleSafeAck(*IncomingSafeAck)
	HandleSafeReadRequest(*IncomingSafeReadRequest)
	HandleError(*IncomingInvalidMessage)
}

type PluggableHandler struct {
	AckHandler         func(*IncomingSafeAck)
	ReadRequestHandler func(*IncomingSafeReadRequest)
	ErrorHandler       func(*IncomingInvalidMessage)
}

func (h *PluggableHandler) HandleSafeAck(ack *IncomingSafeAck) {
	h.AckHandler(ack)
}

func (h *PluggableHandler) HandleSafeReadRequest(read *IncomingSafeReadRequest) {
	h.ReadRequestHandler(read)
}

func (h *PluggableHandler) HandleError(invalid *IncomingInvalidMessage) {
	h.ErrorHandler(invalid)
}
