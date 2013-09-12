package serverconfig

import (
	"io"
	"net"
	"os"
	"path"
	"time"

	"github.com/mark-rushakoff/go_tftpd/readsession"
	"github.com/mark-rushakoff/go_tftpd/readsessioncollection"
	"github.com/mark-rushakoff/go_tftpd/responseagent"
	"github.com/mark-rushakoff/go_tftpd/safepacketprovider"
	"github.com/mark-rushakoff/go_tftpd/sessioncreator"
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
	provider := safepacketprovider.NewSafePacketProvider(c.PacketConn)

	go func() {
		for {
			provider.Read()
		}
	}()

	sessions := readsessioncollection.NewReadSessionCollection()
	sessionCreator := sessioncreator.NewSessionCreator(sessions, readerFromFilename, c.outgoingHandlerFromAddr(), c.DefaultTimeout, c.TryLimit)
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

func (c *ServerConfig) outgoingHandlerFromAddr() sessioncreator.OutgoingHandlerFromAddr {
	return func(addr net.Addr) readsession.OutgoingHandler {
		return responseagent.NewResponseAgent(c.PacketConn, addr)
	}
}
