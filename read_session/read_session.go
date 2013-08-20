package read_session

import (
	"net"

	"github.com/mark-rushakoff/go_tftpd/response_agent"
	"github.com/mark-rushakoff/go_tftpd/safe_packets"
)

type ReadSession struct {
	conn          net.PacketConn
	addr          net.Addr
	responseAgent response_agent.ResponderAgent
}

func NewReadSession(conn net.PacketConn, addr net.Addr, responseAgent response_agent.ResponderAgent, readRequest *safe_packets.SafeReadRequest) *ReadSession {
	session := &ReadSession{
		conn,
		addr,
		responseAgent,
	}

	responseAgent.SendAck(safe_packets.NewSafeAck(0))

	return session
}
