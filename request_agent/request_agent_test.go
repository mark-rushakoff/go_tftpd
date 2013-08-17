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
	const blockNum uint16 = 3456

	conn := &helpers.MockPacketConn{}
	wasCalledOnce := false
	conn.ReadFromFunc = func(b []byte) (int, net.Addr, error) {
		if wasCalledOnce {
			// block forever
			select {}
		}
		var data = []interface{}{
			uint16(messages.ErrorOpcode),
			uint16(messages.FileNotFound),
			string("lol"),
			byte(0),
		}

		buf := new(bytes.Buffer)
		for _, v := range data {
			str, isString := v.(string)
			if isString {
				_, err := buf.WriteString(str)
				if err != nil {
					t.Errorf("Failed to write string to buffer")
				}
			} else {
				err := binary.Write(buf, binary.BigEndian, v)
				if err != nil {
					t.Errorf("Failed to write data to buffer")
				}
			}
		}
		n := copy(b, buf.Bytes())
		return n, nil, nil
	}

	agent := NewRequestAgent(conn)
	go agent.Read()

	select {
	case errorPacket := <-agent.Error:
		expectedCode := messages.FileNotFound
		if errorPacket.Code != expectedCode {
			t.Errorf("Received code %v, expected %v", errorPacket.Code, expectedCode)
		}

		expectedMessage := "lol"
		if errorPacket.Message != expectedMessage {
			t.Errorf("Received message %v, expected %v", errorPacket.Message, expectedMessage)
		}
	case <-time.After(timeoutMs * time.Millisecond):
		t.Errorf("Did not receive Error in time")
	}
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
	const blockNum uint16 = 9876

	conn := &helpers.MockPacketConn{}
	wasCalledOnce := false
	conn.ReadFromFunc = func(b []byte) (int, net.Addr, error) {
		if wasCalledOnce {
			// block forever
			select {}
		}
		var data = []interface{}{
			uint16(messages.ReadOpcode),
			string("/foo/bar"),
			byte(0),
			string("netascii"),
			byte(0),
		}

		buf := new(bytes.Buffer)
		for _, v := range data {
			str, isString := v.(string)
			if isString {
				_, err := buf.WriteString(str)
				if err != nil {
					t.Errorf("Failed to write string to buffer")
				}
			} else {
				err := binary.Write(buf, binary.BigEndian, v)
				if err != nil {
					t.Errorf("Failed to write data to buffer")
				}
			}
		}
		n := copy(b, buf.Bytes())
		return n, nil, nil
	}

	agent := NewRequestAgent(conn)
	go agent.Read()

	select {
	case readPacket := <-agent.ReadRequest:
		expectedFilename := "/foo/bar"
		if readPacket.Filename != expectedFilename {
			t.Errorf("Received name %v, expected %v", readPacket.Filename, expectedFilename)
		}

		expectedMode := messages.NetAscii
		if readPacket.Mode != expectedMode {
			t.Errorf("Received mode %v, expected %v", readPacket.Mode, expectedMode)
		}
	case <-time.After(timeoutMs * time.Millisecond):
		t.Errorf("Did not receive Read in time")
	}
}

func TestWriteRequestPacketCausesWriteRequest(t *testing.T) {
	t.Skipf("Pending")
}
