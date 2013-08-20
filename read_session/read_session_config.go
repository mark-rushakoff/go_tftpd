package read_session

import (
	"io"

	"github.com/mark-rushakoff/go_tftpd/response_agent"
)

type ReadSessionConfig struct {
	ResponseAgent response_agent.ResponderAgent
	Reader        io.Reader

	BlockSize uint16
}
