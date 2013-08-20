package read_session

import (
	"testing"

	"github.com/mark-rushakoff/go_tftpd/response_agent"
	"github.com/mark-rushakoff/go_tftpd/safe_packets"
	"github.com/mark-rushakoff/go_tftpd/test_helpers"
)

func TestNewReadSession(t *testing.T) {
	responseAgent := response_agent.MakeMockResponseAgent()
	readRequest := safe_packets.NewSafeReadRequest("foo", safe_packets.NetAscii)
	conn := &test_helpers.MockPacketConn{}
	addr := &test_helpers.MockAddr{}

	NewReadSession(conn, addr, responseAgent, readRequest)

	if responseAgent.TotalMessagesSent() != 1 {
		t.Fatalf("Expected 1 message sent but %v messages were sent", responseAgent.TotalMessagesSent())
	}

	sentAck := responseAgent.MostRecentAck()
	actualBlockNumber := sentAck.BlockNumber
	expectedBlockNumber := uint16(0)
	if actualBlockNumber != expectedBlockNumber {
		t.Errorf("Expected ReadSession to ack with block number %v, received %v", expectedBlockNumber, actualBlockNumber)
	}
}
