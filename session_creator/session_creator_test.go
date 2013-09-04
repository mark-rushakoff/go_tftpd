package session_creator

import (
	"io"
	"net"
	"testing"
	"time"

	"github.com/mark-rushakoff/go_tftpd/read_session"
	"github.com/mark-rushakoff/go_tftpd/read_session_collection"
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

	readSessions := read_session_collection.NewReadSessionCollection()
	reader := make(chan []byte)
	outgoing := make(chan *safe_packets.SafeData, 1)
	sessionCreator := NewSessionCreator(readSessions, readerFactory(reader), outgoingFactory(outgoing))

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
	}

	_, found := readSessions.Fetch(fakeAddr)
	if !found {
		t.Fatalf("SessionCreator did not expose the session")
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

	return func(filename string) io.Reader {
		return reader
	}
}

type channelNotifier struct {
	Out chan<- *safe_packets.SafeData
}

func (n *channelNotifier) SendData(data *safe_packets.SafeData) {
	n.Out <- data
}

func outgoingFactory(out chan *safe_packets.SafeData) OutgoingHandlerFromAddr {
	return func(net.Addr) read_session.OutgoingHandler {
		return &channelNotifier{
			Out: out,
		}
	}
}
