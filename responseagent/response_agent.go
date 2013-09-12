package responseagent

import (
	"net"

	"github.com/mark-rushakoff/go_tftpd/safepackets"
)

type ResponseAgent struct {
	conn       net.PacketConn
	clientAddr net.Addr
}

// Takes safe packets and serializes them and sends them out on the associated connection.
type ResponderAgent interface {
	SendAck(ack *safepackets.SafeAck)
	SendError(e *safepackets.SafeError)
	SendData(data *safepackets.SafeData)
}

func NewResponseAgent(conn net.PacketConn, clientAddr net.Addr) *ResponseAgent {
	return &ResponseAgent{
		conn:       conn,
		clientAddr: clientAddr,
	}
}

func (a *ResponseAgent) SendAck(ack *safepackets.SafeAck) {
	a.conn.WriteTo(ack.Bytes(), a.clientAddr)
}

func (a *ResponseAgent) SendError(e *safepackets.SafeError) {
	a.conn.WriteTo(e.Bytes(), a.clientAddr)
}

func (a *ResponseAgent) SendData(data *safepackets.SafeData) {
	a.conn.WriteTo(data.Bytes(), a.clientAddr)
}
