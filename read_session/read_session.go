package read_session

import (
	"io"
	"net"

	"github.com/mark-rushakoff/go_tftpd/response_agent"
	"github.com/mark-rushakoff/go_tftpd/safe_packets"
)

type ReadSession struct {
	conn          net.PacketConn
	addr          net.Addr
	responseAgent response_agent.ResponderAgent
}

func NewReadSession(conn net.PacketConn, addr net.Addr, responseAgent response_agent.ResponderAgent, reader io.Reader) *ReadSession {
	session := &ReadSession{
		conn,
		addr,
		responseAgent,
	}

	responseAgent.SendAck(safe_packets.NewSafeAck(0))

	dataBytes := make([]byte, 512)
	bytesRead, _ := reader.Read(dataBytes) // TODO: be defensive
	dataBytes = dataBytes[:bytesRead]
	responseAgent.SendData(safe_packets.NewSafeData(1, dataBytes))

	return session
}
