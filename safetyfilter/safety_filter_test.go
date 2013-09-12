package safetyfilter

import (
	"testing"
	"time"

	"github.com/mark-rushakoff/go_tftpd/packets"
	"github.com/mark-rushakoff/go_tftpd/requestagent"
	"github.com/mark-rushakoff/go_tftpd/safepackets"
	"github.com/mark-rushakoff/go_tftpd/test_helpers"
)

var fakeAddr = test_helpers.MakeMockAddr("fake_network", "a fake addr")

func TestConvertsAcksToSafeAcks(t *testing.T) {
	incomingAcks := make(chan *IncomingSafeAck, 1)
	handler := &PluggableHandler{
		AckHandler: func(ack *IncomingSafeAck) {
			incomingAcks <- ack
		},
	}

	expectedBlockNumber := uint16(500)

	ackPacket := &packets.Ack{
		BlockNumber: expectedBlockNumber,
	}

	ack := &requestagent.IncomingAck{
		Ack:  ackPacket,
		Addr: fakeAddr,
	}

	MakeSafetyFilter(handler).HandleAck(ack)

	select {
	case incomingAck := <-incomingAcks:
		if incomingAck == nil {
			t.Fatalf("Did not receive Ack")
		}

		if incomingAck.Addr.String() != fakeAddr.String() {
			t.Errorf("Received incorrect addr: %v", incomingAck.Addr)
		}

		actualBlockNumber := incomingAck.Ack.BlockNumber
		if actualBlockNumber != expectedBlockNumber {
			t.Errorf("Expected ack with block number %v, but received %v", actualBlockNumber, expectedBlockNumber)
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

	expectedFilename := "foobar"
	modeString := "netascii"
	expectedMode := safepackets.NetAscii

	fakeIncomingReadPacket := &packets.ReadRequest{
		Filename: expectedFilename,
		Mode:     modeString,
	}

	fakeIncomingReadRequest := &requestagent.IncomingReadRequest{
		Read: fakeIncomingReadPacket,
		Addr: fakeAddr,
	}

	MakeSafetyFilter(handler).HandleReadRequest(fakeIncomingReadRequest)

	select {
	case incomingRead := <-incomingReadRequests:
		if incomingRead == nil {
			t.Fatalf("Did not receive Read")
		}

		if incomingRead.Addr.String() != fakeAddr.String() {
			t.Errorf("Received incorrect addr: %v", incomingRead.Addr)
		}

		actualFilename := incomingRead.Read.Filename
		if actualFilename != expectedFilename {
			t.Errorf("Expected Filename '%v', but received '%v'", actualFilename, expectedFilename)
		}

		actualMode := incomingRead.Read.Mode
		if actualMode != expectedMode {
			t.Errorf("Expected Mode '%v', but received '%v'", actualMode, expectedMode)
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

	modeString := "an invalid mode"
	var expectedCode packets.ErrorCode = packets.Undefined
	expectedMessage := "Invalid mode string"

	fakeIncomingReadPacket := &packets.ReadRequest{
		Filename: "foobar",
		Mode:     modeString,
	}

	fakeIncomingReadRequest := &requestagent.IncomingReadRequest{
		Read: fakeIncomingReadPacket,
		Addr: fakeAddr,
	}

	MakeSafetyFilter(handler).HandleReadRequest(fakeIncomingReadRequest)

	select {
	case invalid := <-incomingInvalidMessages:
		if invalid == nil {
			t.Fatalf("Did not receive invalid message")
		}

		if invalid.Addr.String() != fakeAddr.String() {
			t.Errorf("Received incorrect addr: %v", invalid.Addr)
		}

		if invalid.ErrorCode != expectedCode {
			t.Errorf("Received code %v, expected code %v", invalid.ErrorCode, expectedCode)
		}

		if invalid.ErrorMessage != expectedMessage {
			t.Errorf("Received error message '%v', expected message '%v'", invalid.ErrorMessage, expectedMessage)
		}

	case <-time.After(time.Millisecond):
		t.Fatalf("Did not receive read request in time")
	}
}
