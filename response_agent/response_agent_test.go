package response_agent

import (
	"bytes"
	"net"
	"testing"

	"github.com/mark-rushakoff/go_tftpd/test_helpers"
	"github.com/mark-rushakoff/go_tftpd/safe_packets"
)

func TestNewResponseAgent(t *testing.T) {
	agent := NewResponseAgent(nil, nil)
	if agent == nil {
		t.Errorf("Unable to create agent")
	}
}

func TestAckSerializes(t *testing.T) {
	expectedPacketOut := []byte{0, 4, 4, 210}
	agent, conn, addr := buildAgentExpecting(expectedPacketOut)
	ack := safe_packets.NewSafeAck(1234)
	agent.SendAck(ack)

	lastPacketOut, lastAddr, ok := conn.LastPacketOut()
	if !ok {
		t.Errorf("Expected a packet to be sent but no packets were sent")
	}

	if addr != lastAddr {
		t.Errorf("Expected agent to send to addr %v, but sent to %v", addr, lastAddr)
	}

	if !bytes.Equal(expectedPacketOut, lastPacketOut) {
		t.Errorf("Expected outgoing packet %v, received %v", expectedPacketOut, lastPacketOut)
	}
}

func buildAgentExpecting(b []byte) (agent *ResponseAgent, conn *test_helpers.MockPacketConn, addr net.Addr) {
	conn = &test_helpers.MockPacketConn{
		WriteToFunc: func(b []byte, a net.Addr) (int, error) {
			return len(b), nil
		},
	}
	addr = &test_helpers.MockAddr{}
	agent = NewResponseAgent(conn, addr)

	return
}
