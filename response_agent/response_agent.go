package response_agent

import (
	"net"

	"github.com/mark-rushakoff/go_tftpd/safe_packets"
)

type ResponseAgent struct {
	conn       net.PacketConn
	clientAddr net.Addr
}

type ResponderAgent interface {
	SendAck(ack *safe_packets.SafeAck)
	SendErrorPacket(e *safe_packets.SafeError)
	SendData(data *safe_packets.SafeData)
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

func (a *ResponseAgent) SendData(data *safe_packets.SafeData) {
	a.conn.WriteTo(data.Bytes(), a.clientAddr)
}
