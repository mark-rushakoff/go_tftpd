package response_agent

import (
	"bytes"
	"net"
	"testing"

	"github.com/mark-rushakoff/go_tftpd/safe_packets"
	"github.com/mark-rushakoff/go_tftpd/test_helpers"
)

func TestNewResponseAgent(t *testing.T) {
	agent := NewResponseAgent(nil, nil)
	if agent == nil {
		t.Errorf("Unable to create agent")
	}
}

func TestAckSerializes(t *testing.T) {
	expectedPacketOut := []byte{0, 4, 4, 210}
	agent, conn, addr := buildAgentThatWrites(expectedPacketOut)
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

func TestErrorSerializes(t *testing.T) {
	// opcode 5, error code 1, string "File not found"
	expectedPacketOut := []byte{0, 5, 0, 1, 0x46, 0x69, 0x6c, 0x65, 0x20, 0x6e, 0x6f, 0x74, 0x20, 0x66, 0x6f, 0x75, 0x6e, 0x64, 0}
	agent, conn, addr := buildAgentThatWrites(expectedPacketOut)
	errorPacket := safe_packets.NewFileNotFoundError()
	agent.SendErrorPacket(errorPacket)

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

func TestDataSerializes(t *testing.T) {
	// opcode 3, block number 1234, string "foo"
	expectedPacketOut := []byte{0, 3, 4, 210, 102, 111, 111}
	agent, conn, addr := buildAgentThatWrites(expectedPacketOut)
	data := safe_packets.NewSafeData(1234, []byte("foo"))
	agent.SendData(data)

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

func buildAgentThatWrites(b []byte) (agent *ResponseAgent, conn *test_helpers.MockPacketConn, addr net.Addr) {
	conn = &test_helpers.MockPacketConn{
		WriteToFunc: func(b []byte, a net.Addr) (int, error) {
			return len(b), nil
		},
	}
	addr = &test_helpers.MockAddr{}
	agent = NewResponseAgent(conn, addr)

	return
}