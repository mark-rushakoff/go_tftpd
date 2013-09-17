package safetyfilter

import (
	"testing"
	"time"

	"github.com/mark-rushakoff/go_tftpd/packets"
	"github.com/mark-rushakoff/go_tftpd/requestagent"
	"github.com/mark-rushakoff/go_tftpd/safepackets"
	"github.com/mark-rushakoff/go_tftpd/testhelpers"
)

var fakeAddr = testhelpers.MakeMockAddr("fake_network", "a fake addr")

func TestConvertsAcksToSafeAcks(t *testing.T) {
	incomingAcks := make(chan *IncomingSafeAck, 1)
	handler := &PluggableHandler{
		AckHandler: func(ack *IncomingSafeAck) {
			incomingAcks <- ack
		},
	}

	ackPacket := &packets.Ack{
		BlockNumber: 500,
	}

	fakeSafeAck := safepackets.NewSafeAck(500)
	fakeConverter := &safepackets.PluggableConverter{
		FromAckHandler: func(ack *packets.Ack) *safepackets.SafeAck {
			if ack != ackPacket {
				t.Fatalf("fakeConverter called with unexpected argument")
			}

			return fakeSafeAck
		},
	}

	ack := &requestagent.IncomingAck{
		Ack:  ackPacket,
		Addr: fakeAddr,
	}

	MakeSafetyFilter(fakeConverter, handler).HandleAck(ack)

	select {
	case incomingAck := <-incomingAcks:
		if incomingAck.Ack != fakeSafeAck {
			t.Fatalf("SafetyFilter did not use ack provided by converter")
		}

		if incomingAck.Addr != fakeAddr {
			t.Fatalf("SafetyFilter did not use correct addr on incoming ack")
		}
	case <-time.After(time.Millisecond):
		t.Fatalf("Did not receive ack in time")
	}
}

func TestConvertsReadRequestsToSafeReadRequests(t *testing.T) {
	incomingReadRequests := make(chan *IncomingSafeReadRequest, 1)
	handler := &PluggableHandler{
		ReadRequestHandler: func(read *IncomingSafeReadRequest) {
			incomingReadRequests <- read
		},
	}

	fakeIncomingReadPacket := &packets.ReadRequest{
		Filename: "some file",
		Mode:     "netascii",
	}

	fakeSafeRead := safepackets.NewSafeReadRequest("some file", safepackets.NetAscii)
	fakeConverter := &safepackets.PluggableConverter{
		FromReadRequestHandler: func(read *packets.ReadRequest) (*safepackets.SafeReadRequest, *safepackets.ConversionError) {
			if read != fakeIncomingReadPacket {
				t.Fatalf("fakeConverter called with unexpected argument")
			}

			return fakeSafeRead, nil
		},
	}

	fakeIncomingReadRequest := &requestagent.IncomingReadRequest{
		Read: fakeIncomingReadPacket,
		Addr: fakeAddr,
	}
	MakeSafetyFilter(fakeConverter, handler).HandleReadRequest(fakeIncomingReadRequest)

	select {
	case incomingRead := <-incomingReadRequests:
		if incomingRead == nil {
			t.Fatalf("Did not receive Read")
		}

		if incomingRead.Addr.String() != fakeAddr.String() {
			t.Errorf("Received incorrect addr: %v", incomingRead.Addr)
		}

		if incomingRead.Read != fakeSafeRead {
			t.Fatalf("SafetyFilter did not use read provided by fake converter")
		}

	case <-time.After(time.Millisecond):
		t.Fatalf("Did not receive read request in time")
	}
}

func TestRejectsReadRequestWithInvalidMode(t *testing.T) {
	incomingInvalidMessages := make(chan *IncomingInvalidMessage, 1)
	handler := &PluggableHandler{
		ErrorHandler: func(message *IncomingInvalidMessage) {
			incomingInvalidMessages <- message
		},
	}

	fakeIncomingReadPacket := &packets.ReadRequest{
		Filename: "foobar",
		Mode:     "an invalid mode",
	}

	fakeIncomingReadRequest := &requestagent.IncomingReadRequest{
		Read: fakeIncomingReadPacket,
		Addr: fakeAddr,
	}

	fakeConverter := &safepackets.PluggableConverter{
		FromReadRequestHandler: func(read *packets.ReadRequest) (*safepackets.SafeReadRequest, *safepackets.ConversionError) {
			if read != fakeIncomingReadPacket {
				t.Fatalf("fakeConverter called with unexpected argument")
			}

			return nil, safepackets.NewConversionError(packets.Undefined, "Invalid mode string")
		},
	}

	MakeSafetyFilter(fakeConverter, handler).HandleReadRequest(fakeIncomingReadRequest)

	select {
	case invalid := <-incomingInvalidMessages:
		if invalid == nil {
			t.Fatalf("Did not receive invalid message")
		}

		if invalid.Addr.String() != fakeAddr.String() {
			t.Errorf("Received incorrect addr: %v", invalid.Addr)
		}

		if invalid.ErrorCode != packets.Undefined {
			t.Errorf("Received code %v, expected code %v", invalid.ErrorCode, packets.Undefined)
		}

		if invalid.ErrorMessage != "Invalid mode string" {
			t.Errorf("Received error message '%v', expected message '%v'", invalid.ErrorMessage, "Invalid mode string")
		}

	case <-time.After(time.Millisecond):
		t.Fatalf("Did not receive read request in time")
	}
}
