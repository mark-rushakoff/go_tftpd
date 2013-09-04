package session_creator

import (
	"io"
	"net"
	"time"

	"github.com/mark-rushakoff/go_tftpd/read_session"
	"github.com/mark-rushakoff/go_tftpd/read_session_collection"
	"github.com/mark-rushakoff/go_tftpd/safety_filter"
	"github.com/mark-rushakoff/go_tftpd/timeout_controller"
)

type ReaderFromFilename func(filename string) io.Reader
type OutgoingHandlerFromAddr func(net.Addr) read_session.OutgoingHandler

type SessionCreator struct {
	readSessions           *read_session_collection.ReadSessionCollection
	readerFactory          ReaderFromFilename
	outgoingHandlerFactory OutgoingHandlerFromAddr
	timeout                time.Duration
	tryLimit               uint
}

func NewSessionCreator(
	readSessions *read_session_collection.ReadSessionCollection,
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
	sessionConfig := &read_session.Config{
		Reader:    c.readerFactory(r.Read.Filename),
		BlockSize: 512,
	}

	session := read_session.NewReadSession(sessionConfig, c.outgoingHandlerFactory(r.Addr), func() {
		// nothing - should remove from readSessions but needs tests!
	})

	timeoutController := timeout_controller.NewTimeoutController(c.timeout, c.tryLimit, session, func() {
		c.readSessions.Remove(r.Addr)
	})

	c.readSessions.Add(timeoutController, r.Addr)
	go timeoutController.BeginSession()
}
