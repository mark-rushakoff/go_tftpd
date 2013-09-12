package readsession

import (
	"github.com/mark-rushakoff/go_tftpd/safepackets"
)

type OutgoingHandler interface {
	SendData(*safepackets.SafeData)
	SendError(*safepackets.SafeError)
}

type PluggableHandler struct {
	SendDataHandler  func(*safepackets.SafeData)
	SendErrorHandler func(*safepackets.SafeError)
}

func (h *PluggableHandler) SendData(data *safepackets.SafeData) {
	h.SendDataHandler(data)
}

func (h *PluggableHandler) SendError(e *safepackets.SafeError) {
	h.SendErrorHandler(e)
}
