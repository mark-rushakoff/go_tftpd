package request_agent

import (
	"bytes"
	"encoding/binary"
	"net"
	"testing"
	"time"

	"github.com/mark-rushakoff/go_tftpd/helpers"
	"github.com/mark-rushakoff/go_tftpd/packets"
)

const timeoutMs = 100

func TestAcknowledgementPacketCausesAck(t *testing.T) {
	const blockNum uint16 = 1234

	agent := agentWithIncomingPacket(t, []interface{}{
		uint16(packets.AckOpcode),
		uint16(blockNum),
	})

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

	agent := agentWithIncomingPacket(t, []interface{}{
		uint16(packets.ErrorOpcode),
		uint16(packets.FileNotFound),
		string("lol"),
		byte(0),
	})

	select {
	case errorPacket := <-agent.Error:
		expectedCode := packets.FileNotFound
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

	agent := agentWithIncomingPacket(t, []interface{}{
		uint16(packets.DataOpcode),
		uint16(blockNum),
		[]byte{0, 1, 2, 3, 4, 5, 255},
	})

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

	agent := agentWithIncomingPacket(t, []interface{}{
		uint16(packets.ReadOpcode),
		string("/foo/bar"),
		byte(0),
		string("netascii"),
		byte(0),
	})

	select {
	case readPacket := <-agent.ReadRequest:
		expectedFilename := "/foo/bar"
		if readPacket.Filename != expectedFilename {
			t.Errorf("Received name %v, expected %v", readPacket.Filename, expectedFilename)
		}

		expectedMode := "netascii"
		if readPacket.Mode != expectedMode {
			t.Errorf("Received mode %v, expected %v", readPacket.Mode, expectedMode)
		}
	case <-time.After(timeoutMs * time.Millisecond):
		t.Errorf("Did not receive Read in time")
	}
}

func TestWriteRequestPacketCausesWriteRequest(t *testing.T) {
	const blockNum uint16 = 2468

	agent := agentWithIncomingPacket(t, []interface{}{
		uint16(packets.WriteOpcode),
		string("/foo/bar"),
		byte(0),
		string("netascii"),
		byte(0),
	})

	select {
	case writePacket := <-agent.WriteRequest:
		expectedFilename := "/foo/bar"
		if writePacket.Filename != expectedFilename {
			t.Errorf("Received name %v, expected %v", writePacket.Filename, expectedFilename)
		}

		expectedMode := "netascii"
		if writePacket.Mode != expectedMode {
			t.Errorf("Received mode %v, expected %v", writePacket.Mode, expectedMode)
		}
	case <-time.After(timeoutMs * time.Millisecond):
		t.Errorf("Did not receive Read in time")
	}
}

type invalidPacketTestCase struct {
	packet      []byte
	reason      InvalidTransmissionReason
	description string
}

func TestInvalidPacketCausesInvalidTransmission(t *testing.T) {
	testCases := []invalidPacketTestCase{
		{[]byte{}, PacketTooShort, "0-byte packet too short"},
		{[]byte{0}, PacketTooShort, "1-byte packet too short"},
		{[]byte{0, 0}, PacketTooShort, "2-byte packet too short"},
		{[]byte{0, 1, 102, 111, 111}, MissingField, "Read packet with missing filename terminator"},
		{[]byte{0, 1, 102, 111, 111, 0, 98, 97, 114}, MissingField, "Read packet with missing mode terminator"},
		{[]byte{255, 255, 255, 255}, InvalidOpcode, "Invalid opcode"},
	}

	for _, testCase := range testCases {
		invalidPacketCausesInvalidTransmission(t, testCase)
	}
}

func invalidPacketCausesInvalidTransmission(t *testing.T, testCase invalidPacketTestCase) {
	t.Logf(`Running test case "%v"`, testCase.description)
	invalidPacket := testCase.packet
	agent := agentWithIncomingPacket(t, []interface{}{invalidPacket})
	select {
	case invalidTransmission := <-agent.InvalidTransmission:
		expectedTransmission := make([]byte, len(invalidPacket))
		copy(expectedTransmission, invalidPacket)
		if !bytes.Equal(invalidTransmission.Packet, expectedTransmission) {
			t.Errorf("Detected invalid transmission %v, expected %v", invalidTransmission, expectedTransmission)
		}

		actualReason := invalidTransmission.Reason
		if actualReason != testCase.reason {
			t.Errorf("Detected invalid transmission with reason code '%v', expected '%v'", actualReason, testCase.reason)
		}
	case <-time.After(timeoutMs * time.Millisecond):
		t.Errorf("Did not receive invalid transmission in time")
	}
}

func agentWithIncomingPacket(t *testing.T, data []interface{}) *RequestAgent {
	conn := &helpers.MockPacketConn{
		ReadFromFunc: buildReaderFunc(t, data),
	}

	agent := NewRequestAgent(conn)
	go agent.Read()

	return agent
}

func buildReaderFunc(t *testing.T, data []interface{}) func([]byte) (int, net.Addr, error) {
	wasCalledOnce := false
	return func(b []byte) (int, net.Addr, error) {
		if wasCalledOnce {
			// block forever
			select {}
		}
		wasCalledOnce = true

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
}
