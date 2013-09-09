package read_session

import (
	"github.com/mark-rushakoff/go_tftpd/safe_packets"
)

type OutgoingHandler interface {
	SendData(*safe_packets.SafeData)
	SendError(*safe_packets.SafeError)
}

type PluggableHandler struct {
	SendDataHandler  func(*safe_packets.SafeData)
	SendErrorHandler func(*safe_packets.SafeError)
}

func (h *PluggableHandler) SendData(data *safe_packets.SafeData) {
	h.SendDataHandler(data)
}

func (h *PluggableHandler) SendError(e *safe_packets.SafeError) {
	h.SendErrorHandler(e)
}
