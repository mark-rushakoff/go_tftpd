package requestagent

import (
	"bytes"
	"encoding/binary"
	"net"
	"testing"
	"time"

	"github.com/mark-rushakoff/go_tftpd/packets"
	"github.com/mark-rushakoff/go_tftpd/test_helpers"
)

var fakeAddr = test_helpers.MakeMockAddr("fake_network", "a fake addr")

func TestAcknowledgementPacketCausesAck(t *testing.T) {
	const blockNum uint16 = 1234

	incomingAcks := make(chan *IncomingAck, 1)
	handler := &PluggableHandler{
		AckHandler: func(ack *IncomingAck) {
			incomingAcks <- ack
		},
	}
	agentWithIncomingPacket(t, handler, []interface{}{
		uint16(packets.AckOpcode),
		uint16(blockNum),
	}).Read()

	select {
	case incomingAck := <-incomingAcks:
		if incomingAck == nil {
			t.Fatalf("Did not receive Ack")
		}

		if incomingAck.Addr.String() != fakeAddr.String() {
			t.Errorf("Received incorrect addr: %v", incomingAck.Addr)
		}

		ack := incomingAck.Ack
		if ack.BlockNumber != blockNum {
			t.Errorf("Received block number %v, expected %v", ack.BlockNumber, blockNum)
		}

	case <-time.After(time.Millisecond):
		t.Errorf("Did not receive Ack in time")
	}
}

func TestErrorPacketCausesError(t *testing.T) {
	const blockNum uint16 = 3456

	incomingErrors := make(chan *IncomingError, 1)
	handler := &PluggableHandler{
		ErrorHandler: func(e *IncomingError) {
			incomingErrors <- e
		},
	}
	agentWithIncomingPacket(t, handler, []interface{}{
		uint16(packets.ErrorOpcode),
		uint16(packets.FileNotFound),
		string("lol"),
		byte(0),
	}).Read()

	select {
	case incomingError := <-incomingErrors:
		if incomingError == nil {
			t.Fatalf("Did not receive Error")
		}

		if incomingError.Addr.String() != fakeAddr.String() {
			t.Errorf("Received incorrect addr: %v", incomingError.Addr)
		}

		errorPacket := incomingError.Error
		expectedCode := packets.FileNotFound
		if errorPacket.Code != expectedCode {
			t.Errorf("Received code %v, expected %v", errorPacket.Code, expectedCode)
		}

		expectedMessage := "lol"
		if errorPacket.Message != expectedMessage {
			t.Errorf("Received message %v, expected %v", errorPacket.Message, expectedMessage)
		}
	case <-time.After(time.Millisecond):
		t.Errorf("Did not receive Error in time")
	}
}

func TestDataPacketCausesData(t *testing.T) {
	const blockNum uint16 = 2345

	incomingData := make(chan *IncomingData, 1)
	handler := &PluggableHandler{
		DataHandler: func(data *IncomingData) {
			incomingData <- data
		},
	}
	agentWithIncomingPacket(t, handler, []interface{}{
		uint16(packets.DataOpcode),
		uint16(blockNum),
		[]byte{0, 1, 2, 3, 4, 5, 255},
	}).Read()

	select {
	case incomingDatum := <-incomingData:
		if incomingDatum == nil {
			t.Fatalf("Did not receive Data")
		}

		if incomingDatum.Addr.String() != fakeAddr.String() {
			t.Errorf("Received incorrect addr: %v", incomingDatum.Addr)
		}

		data := incomingDatum.Data
		if data.BlockNumber != blockNum {
			t.Errorf("Received block number %v, expected %v", data.BlockNumber, blockNum)
		}

		expectedData := []byte{0, 1, 2, 3, 4, 5, 255}
		if !bytes.Equal(data.Data, expectedData) {
			t.Errorf("Received data %v, expected %v", data.Data, expectedData)
		}
	case <-time.After(time.Millisecond):
		t.Errorf("Did not receive Data in time")
	}
}

func TestReadRequestPacketCausesReadRequest(t *testing.T) {
	const blockNum uint16 = 9876

	incomingReadRequests := make(chan *IncomingReadRequest, 1)
	handler := &PluggableHandler{
		ReadRequestHandler: func(read *IncomingReadRequest) {
			incomingReadRequests <- read
		},
	}
	agentWithIncomingPacket(t, handler, []interface{}{
		uint16(packets.ReadOpcode),
		string("/foo/bar"),
		byte(0),
		string("netascii"),
		byte(0),
	}).Read()

	select {
	case incomingRead := <-incomingReadRequests:
		if incomingRead == nil {
			t.Fatalf("Did not receive Read")
		}

		if incomingRead.Addr.String() != fakeAddr.String() {
			t.Errorf("Received incorrect addr: %v", incomingRead.Addr)
		}

		readPacket := incomingRead.Read
		expectedFilename := "/foo/bar"
		if readPacket.Filename != expectedFilename {
			t.Errorf("Received name %v, expected %v", readPacket.Filename, expectedFilename)
		}

		expectedMode := "netascii"
		if readPacket.Mode != expectedMode {
			t.Errorf("Received mode %v, expected %v", readPacket.Mode, expectedMode)
		}
	case <-time.After(time.Millisecond):
		t.Errorf("Did not receive Read in time")
	}
}

func TestReadRequestPacketWithOptionsCausesReadRequest(t *testing.T) {
	const blockNum uint16 = 9876

	incomingReadRequests := make(chan *IncomingReadRequest, 1)
	handler := &PluggableHandler{
		ReadRequestHandler: func(read *IncomingReadRequest) {
			incomingReadRequests <- read
		},
	}
	agentWithIncomingPacket(t, handler, []interface{}{
		uint16(packets.ReadOpcode),
		string("/foo/bar"),
		byte(0),
		string("netascii"),
		byte(0),
		string("key1"),
		byte(0),
		string("value1"),
		byte(0),
		string("key2"),
		byte(0),
		string("value2"),
		byte(0),
	}).Read()

	select {
	case incomingRead := <-incomingReadRequests:
		if incomingRead == nil {
			t.Fatalf("Did not receive Read")
		}

		if incomingRead.Addr.String() != fakeAddr.String() {
			t.Errorf("Received incorrect addr: %v", incomingRead.Addr)
		}

		readPacket := incomingRead.Read
		expectedFilename := "/foo/bar"
		if readPacket.Filename != expectedFilename {
			t.Errorf("Received name %v, expected %v", readPacket.Filename, expectedFilename)
		}

		expectedMode := "netascii"
		if readPacket.Mode != expectedMode {
			t.Errorf("Received mode %v, expected %v", readPacket.Mode, expectedMode)
		}

		if len(readPacket.Options) != 2 {
			t.Errorf("Received %v option(s), expected 2", len(readPacket.Options))
		}
		if readPacket.Options["key1"] != "value1" {
			t.Errorf("Expected option key1:value1")
		}
		if readPacket.Options["key2"] != "value2" {
			t.Errorf("Expected option key2:value2")
		}
	case <-time.After(time.Millisecond):
		t.Errorf("Did not receive Read in time")
	}
}

func TestWriteRequestPacketCausesWriteRequest(t *testing.T) {
	const blockNum uint16 = 2468

	incomingWriteRequests := make(chan *IncomingWriteRequest, 1)
	handler := &PluggableHandler{
		WriteRequestHandler: func(write *IncomingWriteRequest) {
			incomingWriteRequests <- write
		},
	}
	agentWithIncomingPacket(t, handler, []interface{}{
		uint16(packets.WriteOpcode),
		string("/foo/bar"),
		byte(0),
		string("netascii"),
		byte(0),
	}).Read()

	select {
	case incomingWrite := <-incomingWriteRequests:
		if incomingWrite == nil {
			t.Fatalf("Did not receive Read")
		}

		if incomingWrite.Addr.String() != fakeAddr.String() {
			t.Errorf("Received incorrect addr: %v", incomingWrite.Addr)
		}
		writePacket := incomingWrite.Write
		expectedFilename := "/foo/bar"
		if writePacket.Filename != expectedFilename {
			t.Errorf("Received name %v, expected %v", writePacket.Filename, expectedFilename)
		}

		expectedMode := "netascii"
		if writePacket.Mode != expectedMode {
			t.Errorf("Received mode %v, expected %v", writePacket.Mode, expectedMode)
		}
	case <-time.After(time.Millisecond):
		t.Errorf("Did not receive Read in time")
	}
}

func TestWriteRequestPacketWithOptionsCausesWriteRequest(t *testing.T) {
	const blockNum uint16 = 2468

	incomingWriteRequests := make(chan *IncomingWriteRequest, 1)
	handler := &PluggableHandler{
		WriteRequestHandler: func(write *IncomingWriteRequest) {
			incomingWriteRequests <- write
		},
	}
	agentWithIncomingPacket(t, handler, []interface{}{
		uint16(packets.WriteOpcode),
		string("/foo/bar"),
		byte(0),
		string("netascii"),
		byte(0),
		string("key1"),
		byte(0),
		string("value1"),
		byte(0),
		string("key2"),
		byte(0),
		string("value2"),
		byte(0),
	}).Read()

	select {
	case incomingWrite := <-incomingWriteRequests:
		if incomingWrite == nil {
			t.Fatalf("Did not receive Read")
		}

		if incomingWrite.Addr.String() != fakeAddr.String() {
			t.Errorf("Received incorrect addr: %v", incomingWrite.Addr)
		}
		writePacket := incomingWrite.Write
		expectedFilename := "/foo/bar"
		if writePacket.Filename != expectedFilename {
			t.Errorf("Received name %v, expected %v", writePacket.Filename, expectedFilename)
		}

		expectedMode := "netascii"
		if writePacket.Mode != expectedMode {
			t.Errorf("Received mode %v, expected %v", writePacket.Mode, expectedMode)
		}

		if len(writePacket.Options) != 2 {
			t.Errorf("Received %v option(s), expected 2", len(writePacket.Options))
		}
		if writePacket.Options["key1"] != "value1" {
			t.Errorf("Expected option key1:value1")
		}
		if writePacket.Options["key2"] != "value2" {
			t.Errorf("Expected option key2:value2")
		}
	case <-time.After(time.Millisecond):
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
		{[]byte{0, 3, 1}, PacketTooShort, "Data with incomplete block number"},
		{[]byte{0, 4, 1}, PacketTooShort, "Ack with incomplete block number"},
		{[]byte{0, 4, 1, 2, 3}, PacketTooLong, "Ack with extra data"},
		{[]byte{0, 5, 1}, PacketTooShort, "Error with incomplete block number"},
		{[]byte{0, 5, 1, 1, 32, 0, 5}, PacketTooLong, "Error with data after message"},
		{[]byte{0, 1, 102, 111, 111}, MissingField, "Read packet with missing filename terminator"},
		{[]byte{0, 1, 102, 111, 111, 0, 98, 97, 114}, MissingField, "Read packet with missing mode terminator"},
		{[]byte{0, 1, 102, 111, 111, 0, 98, 97, 114, 0, 5}, OptionsMalformed, "Option key missing terminator"},
		{[]byte{0, 1, 102, 111, 111, 0, 98, 97, 114, 0, 98, 0, 98}, OptionsMalformed, "Option value missing terminator"},
		{[]byte{255, 255, 255, 255}, InvalidOpcode, "Invalid opcode"},
	}

	for _, testCase := range testCases {
		invalidPacketCausesInvalidTransmission(t, testCase)
	}
}

func invalidPacketCausesInvalidTransmission(t *testing.T, testCase invalidPacketTestCase) {
	t.Logf(`Running test case "%v"`, testCase.description)
	invalidTransmissions := make(chan *InvalidTransmission, 1)
	handler := &PluggableHandler{
		InvalidTransmissionHandler: func(invalid *InvalidTransmission) {
			invalidTransmissions <- invalid
		},
	}
	invalidPacket := testCase.packet
	agentWithIncomingPacket(t, handler, []interface{}{invalidPacket}).Read()

	select {
	case invalidTransmission := <-invalidTransmissions:
		if invalidTransmission == nil {
			t.Fatalf("Did not receive InvalidTransmission")
		}

		if invalidTransmission.Addr.String() != fakeAddr.String() {
			t.Errorf("Received incorrect addr: %v", invalidTransmission.Addr)
		}
		expectedTransmission := make([]byte, len(invalidPacket))
		copy(expectedTransmission, invalidPacket)
		if !bytes.Equal(invalidTransmission.Packet, expectedTransmission) {
			t.Errorf("Detected invalid transmission %v, expected %v", invalidTransmission, expectedTransmission)
		}

		actualReason := invalidTransmission.Reason
		if actualReason != testCase.reason {
			t.Errorf("Detected invalid transmission with reason code '%v', expected '%v'", actualReason, testCase.reason)
		}
	case <-time.After(time.Millisecond):
		t.Errorf("Did not receive invalid transmission in time")
	}
}

func agentWithIncomingPacket(t *testing.T, handler RequestHandler, data []interface{}) *RequestAgent {
	conn := &test_helpers.MockPacketConn{
		ReadFromFunc: buildReaderFunc(t, data),
	}

	agent := NewRequestAgent(conn, handler)

	return agent
}

func buildReaderFunc(t *testing.T, data []interface{}) func([]byte) (int, net.Addr, error) {
	wasCalledOnce := false
	return func(b []byte) (int, net.Addr, error) {
		if wasCalledOnce {
			t.Fatalf("Called fake reader more than once")
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
		return n, fakeAddr, nil
	}
}
