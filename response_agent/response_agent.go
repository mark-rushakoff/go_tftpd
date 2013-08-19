package response_agent

import (
	"net"

	"github.com/mark-rushakoff/go_tftpd/safe_packets"
)

type ResponseAgent struct {
	conn       net.PacketConn
	clientAddr net.Addr
}

func NewResponseAgent(conn net.PacketConn, clientAddr net.Addr) *ResponseAgent {
	return &ResponseAgent{
		conn:       conn,
		clientAddr: clientAddr,
	}
}

func (a *ResponseAgent) SendAck(ack *safe_packets.SafeAck) {
	a.conn.WriteTo(ack.Bytes(), a.clientAddr)
}

func (a *ResponseAgent) SendErrorPacket(e *safe_packets.SafeError) {
	a.conn.WriteTo(e.Bytes(), a.clientAddr)
}
