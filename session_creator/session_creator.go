package session_creator

import (
	"io"
	"net"

	"github.com/mark-rushakoff/go_tftpd/read_session"
	"github.com/mark-rushakoff/go_tftpd/read_session_collection"
	"github.com/mark-rushakoff/go_tftpd/safety_filter"
)

type ReaderFromFilename func(filename string) io.Reader
type OutgoingHandlerFromAddr func(net.Addr) read_session.OutgoingHandler

type SessionCreator struct {
	readSessions           *read_session_collection.ReadSessionCollection
	readerFactory          ReaderFromFilename
	outgoingHandlerFactory OutgoingHandlerFromAddr
}

func NewSessionCreator(
	readSessions *read_session_collection.ReadSessionCollection,
	readerFactory ReaderFromFilename,
	outgoingHandlerFactory OutgoingHandlerFromAddr,
) *SessionCreator {
	return &SessionCreator{
		readSessions:           readSessions,
		readerFactory:          readerFactory,
		outgoingHandlerFactory: outgoingHandlerFactory,
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

	c.readSessions.Add(session, r.Addr)
	go session.Begin()
}
