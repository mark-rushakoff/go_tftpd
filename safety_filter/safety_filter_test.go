package safety_filter

import (
	"testing"
	"time"

	"github.com/mark-rushakoff/go_tftpd/packets"
	"github.com/mark-rushakoff/go_tftpd/request_agent"
	"github.com/mark-rushakoff/go_tftpd/safe_packets"
)

func TestConvertsAcksToSafeAcks(t *testing.T) {
	requestAgent := request_agent.NewRequestAgent(nil)
	safetyFilter := MakeSafetyFilter(requestAgent)

	expectedBlockNumber := uint16(500)

	go func() {
		outgoingAck := &packets.Ack{
			BlockNumber: expectedBlockNumber,
		}

		outgoing := &request_agent.IncomingAck{
			Ack: outgoingAck,
		}

		requestAgent.Ack <- outgoing
	}()

	select {
	case incomingAck := <-safetyFilter.IncomingAck:
		actualBlockNumber := incomingAck.Ack.BlockNumber
		if actualBlockNumber != expectedBlockNumber {
			t.Errorf("Expected ack with block number %v, but received %v", actualBlockNumber, expectedBlockNumber)
		}
	case <-time.After(20 * time.Millisecond):
		t.Fatalf("Did not receive ack in time")
	}
}

func TestConvertsReadRequestsToSafeReadRequests(t *testing.T) {
	expectedFilename := "foobar"
	modeString := "netascii"
	expectedMode := safe_packets.NetAscii

	requestAgent := request_agent.NewRequestAgent(nil)
	safetyFilter := MakeSafetyFilter(requestAgent)

	go func() {
		fakeIncomingRead := &packets.ReadRequest{
			Filename: expectedFilename,
			Mode:     modeString,
		}

		fakeIncoming := &request_agent.IncomingReadRequest{
			Read: fakeIncomingRead,
		}

		requestAgent.ReadRequest <- fakeIncoming
	}()

	select {
	case incomingRead := <-safetyFilter.IncomingRead:
		actualFilename := incomingRead.Read.Filename
		if actualFilename != expectedFilename {
			t.Errorf("Expected Filename '%v', but received '%v'", actualFilename, expectedFilename)
		}

		actualMode := incomingRead.Read.Mode
		if actualMode != expectedMode {
			t.Errorf("Expected Mode '%v', but received '%v'", actualMode, expectedMode)
		}

	case <-time.After(20 * time.Millisecond):
		t.Fatalf("Did not receive read request in time")
	}
}
