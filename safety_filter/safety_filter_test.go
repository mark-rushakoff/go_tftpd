package safety_filter

import (
	"testing"
	"time"

	"github.com/mark-rushakoff/go_tftpd/packets"
	"github.com/mark-rushakoff/go_tftpd/request_agent"
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

func TestConvertsReadsToSafeReads(t *testing.T) {
}
