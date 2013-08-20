package read_session

import (
	"io"
	"net"

	"github.com/mark-rushakoff/go_tftpd/response_agent"
)

type ReadSessionConfig struct {
	Conn          net.PacketConn
	Addr          net.Addr
	ResponseAgent response_agent.ResponderAgent
	Reader        io.Reader

	BlockSize uint16
}
