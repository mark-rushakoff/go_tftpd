package read_session

import (
	"io"

	"github.com/mark-rushakoff/go_tftpd/response_agent"
	"github.com/mark-rushakoff/go_tftpd/timeout_controller"
)

type ReadSessionConfig struct {
	ResponseAgent     response_agent.ResponderAgent
	Reader            io.Reader
	TimeoutController timeout_controller.TimeoutController

	BlockSize uint16
}
