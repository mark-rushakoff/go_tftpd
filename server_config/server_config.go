package server_config

import (
	"net"
	"os"
	"path"
	"time"

	"github.com/mark-rushakoff/go_tftpd/read_session"
	"github.com/mark-rushakoff/go_tftpd/response_agent"
	"github.com/mark-rushakoff/go_tftpd/safe_packet_provider"
	"github.com/mark-rushakoff/go_tftpd/session_manager"
	"github.com/mark-rushakoff/go_tftpd/timeout_controller"
)

type ServerConfig struct {
	// The PacketConn to use for incoming and outgoing messages
	PacketConn net.PacketConn

	// What block size to use when sending data packets
	BlockSize uint16

	// How long to wait until retrying to send a packet
	Timeout time.Duration

	// How many tries to use when sending a packet until giving up
	TryLimit uint
}

func (c *ServerConfig) Serve() {
	provider := safe_packet_provider.NewSafePacketProvider(c.PacketConn)

	go func() {
		for {
			provider.Read()
		}
	}()

	sessionManager := session_manager.NewSessionManager(c.makeReadSessionFactory())

	for {
		select {
		case r := <-provider.IncomingSafeReadRequest():
			go sessionManager.MakeReadSessionFromIncomingRequest(r)
		case ack := <-provider.IncomingSafeAck():
			sessionManager.SendAckToReadSession(ack)
		}
	}
}

func (c *ServerConfig) makeReadSessionFactory() read_session.ReadSessionFactory {
	return func(filename string, clientAddr net.Addr) *read_session.ReadSession {
		workingDir, err := os.Getwd()
		if err != nil {
			panic(err.Error())
		}

		// serve files as though the filesystem root is the working directory
		file, err := os.Open(path.Join(workingDir, path.Clean(filename)))
		if err != nil {
			panic(err.Error())
		}
		config := &read_session.ReadSessionConfig{
			ResponseAgent:     response_agent.NewResponseAgent(c.PacketConn, clientAddr),
			Reader:            file,
			TimeoutController: timeout_controller.NewTimeoutController(c.Timeout, c.TryLimit),

			BlockSize: c.BlockSize,
		}
		return read_session.NewReadSession(config)
	}
}
