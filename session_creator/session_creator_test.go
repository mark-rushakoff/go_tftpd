package session_creator

import (
	"errors"
	"io"
	"net"
	"testing"
	"time"

	"github.com/mark-rushakoff/go_tftpd/packets"
	"github.com/mark-rushakoff/go_tftpd/readsession"
	"github.com/mark-rushakoff/go_tftpd/readsessioncollection"
	"github.com/mark-rushakoff/go_tftpd/safe_packets"
	"github.com/mark-rushakoff/go_tftpd/safety_filter"
	"github.com/mark-rushakoff/go_tftpd/test_helpers"
)

var fakeAddr = test_helpers.MakeMockAddr("fake_network", "a")

func TestCreateAddsNewSessionToCollection(t *testing.T) {
	readRequest := &safety_filter.IncomingSafeReadRequest{
		Read: safe_packets.NewSafeReadRequest("foobar", safe_packets.NetAscii),
		Addr: fakeAddr,
	}

	readSessions := readsessioncollection.NewReadSessionCollection()
	reader := make(chan []byte)
	outgoing := make(chan *safe_packets.SafeData, 1)
	sessionCreator := NewSessionCreator(
		readSessions,
		readerFactory(reader),
		outgoingFactory(outgoing, nil),
		2*time.Millisecond,
		2,
	)

	sessionCreator.Create(readRequest)

	select {
	case reader <- []byte("foobar"):
		// ok
	case <-time.After(time.Millisecond):
		t.Fatalf("Create did not call begin on the session (because the session did not use the reader)")
	}

	select {
	case data := <-outgoing:
		expected := safe_packets.NewSafeData(1, []byte("foobar"))
		if !data.Equals(expected) {
			t.Fatalf("Session sent wrong data packet: got %v, expected %v", data.Bytes(), expected)
		}
	default:
		t.Fatalf("Session did not send data during BeginSession")
	}

	_, found := readSessions.Fetch(fakeAddr)
	if !found {
		t.Fatalf("SessionCreator did not expose the session")
	}

	select {
	case <-outgoing:
		// don't care
	case <-time.After(3 * time.Millisecond):
		t.Fatalf("Should have timed out and re-sent data")
	}

	// let timeout elapse once more
	time.Sleep(3 * time.Millisecond)

	_, found = readSessions.Fetch(fakeAddr)
	if found {
		t.Fatalf("Should have timed out and removed itself from collection")
	}
}

func TestSuccessfulFinishRemovesSessionFromCollection(t *testing.T) {
	readRequest := &safety_filter.IncomingSafeReadRequest{
		Read: safe_packets.NewSafeReadRequest("foobar", safe_packets.NetAscii),
		Addr: fakeAddr,
	}

	readSessions := readsessioncollection.NewReadSessionCollection()
	reader := make(chan []byte)
	outgoing := make(chan *safe_packets.SafeData, 1)
	sessionCreator := NewSessionCreator(
		readSessions,
		readerFactory(reader),
		outgoingFactory(outgoing, nil),
		2*time.Millisecond,
		2,
	)

	sessionCreator.Create(readRequest)

	select {
	case reader <- []byte("foobar"):
		// ok
	case <-time.After(time.Millisecond):
		t.Fatalf("Create did not call begin on the session (because the session did not use the reader)")
	}

	select {
	case data := <-outgoing:
		expected := safe_packets.NewSafeData(1, []byte("foobar"))
		if !data.Equals(expected) {
			t.Fatalf("Session sent wrong data packet: got %v, expected %v", data.Bytes(), expected)
		}
	default:
		t.Fatalf("Session did not send data during BeginSession")
	}

	session, found := readSessions.Fetch(fakeAddr)
	if !found {
		t.Fatalf("SessionCreator did not expose the session")
	}

	session.HandleAck(safe_packets.NewSafeAck(1))
	_, found = readSessions.Fetch(fakeAddr)
	if found {
		t.Fatalf("Finished session was not removed from collection")
	}
}

func TestErrorCreatingReaderCausesErrorMessage(t *testing.T) {
	readRequest := &safety_filter.IncomingSafeReadRequest{
		Read: safe_packets.NewSafeReadRequest("foobar", safe_packets.NetAscii),
		Addr: fakeAddr,
	}

	err := errors.New("something about foobar")
	readSessions := readsessioncollection.NewReadSessionCollection()
	errors := make(chan *safe_packets.SafeError, 1)
	sessionCreator := NewSessionCreator(
		readSessions,
		errorReaderFactory(err),
		outgoingFactory(nil, errors),
		2*time.Millisecond,
		2,
	)

	sessionCreator.Create(readRequest)
	select {
	case e := <-errors:
		expected := &safe_packets.SafeError{Code: packets.AccessViolation, Message: "something about foobar"}
		if !e.Equals(expected) {
			t.Fatalf("Session sent wrong error packet: got %v, expected %v", e.Bytes(), expected.Bytes())
		}
	default:
		t.Fatalf("Session did not send data during BeginSession")
	}
}

type channelReader struct {
	In <-chan []byte
}

func (r *channelReader) Read(p []byte) (n int, err error) {
	return copy(p, <-r.In), nil
}

func readerFactory(in chan []byte) ReaderFromFilename {
	reader := &channelReader{
		In: in,
	}

	return func(filename string) (io.Reader, error) {
		return reader, nil
	}
}

func errorReaderFactory(err error) ReaderFromFilename {
	return func(string) (io.Reader, error) {
		return nil, err
	}
}

type channelNotifier struct {
	Out chan<- *safe_packets.SafeData
	Err chan<- *safe_packets.SafeError
}

func (n *channelNotifier) SendData(data *safe_packets.SafeData) {
	n.Out <- data
}

func (n *channelNotifier) SendError(err *safe_packets.SafeError) {
	n.Err <- err
}

func outgoingFactory(out chan *safe_packets.SafeData, err chan *safe_packets.SafeError) OutgoingHandlerFromAddr {
	return func(net.Addr) readsession.OutgoingHandler {
		return &channelNotifier{
			Out: out,
			Err: err,
		}
	}
}
