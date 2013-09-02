package read_session

import (
	"github.com/mark-rushakoff/go_tftpd/safe_packets"
)

type OutgoingHandler interface {
	SendData(*safe_packets.SafeData)
}

type PluggableHandler struct {
	SendDataHandler func(*safe_packets.SafeData)
}

func (h *PluggableHandler) SendData(data *safe_packets.SafeData) {
	h.SendDataHandler(data)
}
