package server_config

import (
	"io"
	"net"
	"os"
	"path"
	"time"

	"github.com/mark-rushakoff/go_tftpd/readsession"
	"github.com/mark-rushakoff/go_tftpd/readsessioncollection"
	"github.com/mark-rushakoff/go_tftpd/response_agent"
	"github.com/mark-rushakoff/go_tftpd/safe_packet_provider"
	"github.com/mark-rushakoff/go_tftpd/session_creator"
	"github.com/mark-rushakoff/go_tftpd/session_router"
)

type ServerConfig struct {
	// The PacketConn to use for incoming and outgoing messages
	PacketConn net.PacketConn

	// How long to wait until retrying to send a packet
	DefaultTimeout time.Duration

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

	sessions := readsessioncollection.NewReadSessionCollection()
	sessionCreator := session_creator.NewSessionCreator(sessions, readerFromFilename, c.outgoingHandlerFromAddr(), c.DefaultTimeout, c.TryLimit)
	sessionRouter := session_router.NewSessionRouter(sessions)

	for {
		select {
		case r := <-provider.IncomingSafeReadRequest():
			sessionCreator.Create(r)
		case ack := <-provider.IncomingSafeAck():
			sessionRouter.RouteAck(ack)
		}
	}
}

func readerFromFilename(filename string) (io.Reader, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// serve files as though the filesystem root is the working directory
	file, err := os.Open(path.Join(workingDir, path.Clean(filename)))
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (c *ServerConfig) outgoingHandlerFromAddr() session_creator.OutgoingHandlerFromAddr {
	return func(addr net.Addr) readsession.OutgoingHandler {
		return response_agent.NewResponseAgent(c.PacketConn, addr)
	}
}
