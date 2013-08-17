package request_agent

import (
	"bytes"
	"encoding/binary"
	"net"
	"testing"
	"time"

	"github.com/mark-rushakoff/go_tftpd/helpers"
	"github.com/mark-rushakoff/go_tftpd/messages"
)

const timeoutMs = 100

func TestNewRequestAgent(t *testing.T) {
	agent := NewRequestAgent(nil)
	if agent == nil {
		t.Errorf("Could not create an agent")
	}
}

func TestAcknowledgementPacketCausesAck(t *testing.T) {
	const blockNum uint16 = 1234

	conn := &helpers.MockPacketConn{}
	wasCalledOnce := false
	conn.ReadFromFunc = func(b []byte) (int, net.Addr, error) {
		if wasCalledOnce {
			// block forever
			select {}
		}
		var data = []interface{}{
			uint16(messages.AckOpcode),
			uint16(blockNum),
		}

		buf := new(bytes.Buffer)
		for _, v := range data {
			err := binary.Write(buf, binary.BigEndian, v)
			if err != nil {
				t.Errorf("Failed to write data to buffer")
			}
		}
		n := copy(b, buf.Bytes())
		return n, nil, nil
	}

	agent := NewRequestAgent(conn)
	go agent.Read()

	select {
	case ack := <-agent.Ack:
		if ack.BlockNumber != blockNum {
			t.Errorf("Received block number %v, expected %v", ack.BlockNumber, blockNum)
		}
	case <-time.After(timeoutMs * time.Millisecond):
		t.Errorf("Did not receive Ack in time")
	}
}

func TestErrorPacketCausesError(t *testing.T) {
	t.Skipf("Pending")
}

func TestDataPacketCausesData(t *testing.T) {
	const blockNum uint16 = 2345

	conn := &helpers.MockPacketConn{}
	wasCalledOnce := false
	conn.ReadFromFunc = func(b []byte) (int, net.Addr, error) {
		if wasCalledOnce {
			// block forever
			select {}
		}
		var data = []interface{}{
			uint16(messages.DataOpcode),
			uint16(blockNum),
			[]byte{0, 1, 2, 3, 4, 5, 255},
		}

		buf := new(bytes.Buffer)
		for _, v := range data {
			err := binary.Write(buf, binary.BigEndian, v)
			if err != nil {
				t.Errorf("Failed to write data to buffer")
			}
		}
		n := copy(b, buf.Bytes())
		return n, nil, nil
	}

	agent := NewRequestAgent(conn)
	go agent.Read()

	select {
	case data := <-agent.Data:
		if data.BlockNumber != blockNum {
			t.Errorf("Received block number %v, expected %v", data.BlockNumber, blockNum)
		}

		expectedData := []byte{0, 1, 2, 3, 4, 5, 255}
		if !bytes.Equal(data.Data, expectedData) {
			t.Errorf("Received data %v, expected %v", data.Data, expectedData)
		}
	case <-time.After(timeoutMs * time.Millisecond):
		t.Errorf("Did not receive Data in time")
	}
}

func TestReadRequestPacketCausesReadRequest(t *testing.T) {
	t.Skipf("Pending")
}

func TestWriteRequestPacketCausesWriteRequest(t *testing.T) {
	t.Skipf("Pending")
}
