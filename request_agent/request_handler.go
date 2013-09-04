package request_agent

type RequestHandler interface {
	HandleAck(*IncomingAck)
	HandleError(*IncomingError)
	HandleData(*IncomingData)
	HandleReadRequest(*IncomingReadRequest)
	HandleWriteRequest(*IncomingWriteRequest)
	HandleInvalidTransmission(*InvalidTransmission)
}

type PluggableHandler struct {
	AckHandler                 func(*IncomingAck)
	ErrorHandler               func(*IncomingError)
	DataHandler                func(*IncomingData)
	ReadRequestHandler         func(*IncomingReadRequest)
	WriteRequestHandler        func(*IncomingWriteRequest)
	InvalidTransmissionHandler func(*InvalidTransmission)
}

func (h *PluggableHandler) HandleAck(ack *IncomingAck) {
	h.AckHandler(ack)
}

func (h *PluggableHandler) HandleError(e *IncomingError) {
	h.ErrorHandler(e)
}

func (h *PluggableHandler) HandleData(data *IncomingData) {
	h.DataHandler(data)
}

func (h *PluggableHandler) HandleReadRequest(read *IncomingReadRequest) {
	h.ReadRequestHandler(read)
}

func (h *PluggableHandler) HandleWriteRequest(write *IncomingWriteRequest) {
	h.WriteRequestHandler(write)
}

func (h *PluggableHandler) HandleInvalidTransmission(transmission *InvalidTransmission) {
	h.InvalidTransmissionHandler(transmission)
}
