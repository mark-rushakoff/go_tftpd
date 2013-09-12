package session_creator

import (
	"io"
	"net"
	"time"

	"github.com/mark-rushakoff/go_tftpd/readsession"
	"github.com/mark-rushakoff/go_tftpd/readsessioncollection"
	"github.com/mark-rushakoff/go_tftpd/safe_packets"
	"github.com/mark-rushakoff/go_tftpd/safety_filter"
	"github.com/mark-rushakoff/go_tftpd/timeout_controller"
)

type ReaderFromFilename func(filename string) (io.Reader, error)
type OutgoingHandlerFromAddr func(net.Addr) readsession.OutgoingHandler

type SessionCreator struct {
	readSessions           *readsessioncollection.ReadSessionCollection
	readerFactory          ReaderFromFilename
	outgoingHandlerFactory OutgoingHandlerFromAddr
	timeout                time.Duration
	tryLimit               uint
}

func NewSessionCreator(
	readSessions *readsessioncollection.ReadSessionCollection,
	readerFactory ReaderFromFilename,
	outgoingHandlerFactory OutgoingHandlerFromAddr,
	timeout time.Duration,
	tryLimit uint,
) *SessionCreator {
	return &SessionCreator{
		readSessions:           readSessions,
		readerFactory:          readerFactory,
		outgoingHandlerFactory: outgoingHandlerFactory,
		timeout:                timeout,
		tryLimit:               tryLimit,
	}
}

func (c *SessionCreator) Create(r *safety_filter.IncomingSafeReadRequest) {
	reader, err := c.readerFactory(r.Read.Filename)
	if err != nil {
		handler := c.outgoingHandlerFactory(r.Addr)

		handler.SendError(safe_packets.NewAccessViolationError(err.Error()))
		return
	}

	sessionConfig := &readsession.Config{
		Reader:    reader,
		BlockSize: 512,
	}

	removeSession := func() {
		c.readSessions.Remove(r.Addr)
	}

	session := readsession.NewReadSession(sessionConfig, c.outgoingHandlerFactory(r.Addr), removeSession)

	timeoutController := timeout_controller.NewTimeoutController(c.timeout, c.tryLimit, session, removeSession)

	c.readSessions.Add(timeoutController, r.Addr)
	go timeoutController.BeginSession()
}
