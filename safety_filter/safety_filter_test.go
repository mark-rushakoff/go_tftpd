package safety_filter

import (
	"testing"
	"time"

	"github.com/mark-rushakoff/go_tftpd/packets"
	"github.com/mark-rushakoff/go_tftpd/request_agent"
	"github.com/mark-rushakoff/go_tftpd/safe_packets"
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

	ack := &request_agent.IncomingAck{
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
	expectedMode := safe_packets.NetAscii

	fakeIncomingReadPacket := &packets.ReadRequest{
		Filename: expectedFilename,
		Mode:     modeString,
	}

	fakeIncomingReadRequest := &request_agent.IncomingReadRequest{
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
