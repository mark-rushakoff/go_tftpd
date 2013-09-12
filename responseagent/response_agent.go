package responseagent

import (
	"net"

	"github.com/mark-rushakoff/go_tftpd/safe_packets"
)

type ResponseAgent struct {
	conn       net.PacketConn
	clientAddr net.Addr
}

// Takes safe packets and serializes them and sends them out on the associated connection.
type ResponderAgent interface {
	SendAck(ack *safe_packets.SafeAck)
	SendError(e *safe_packets.SafeError)
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

func (a *ResponseAgent) SendError(e *safe_packets.SafeError) {
	a.conn.WriteTo(e.Bytes(), a.clientAddr)
}

func (a *ResponseAgent) SendData(data *safe_packets.SafeData) {
	a.conn.WriteTo(data.Bytes(), a.clientAddr)
}
