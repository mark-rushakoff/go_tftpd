package read_session

import (
	"bytes"
	"strings"
	"testing"

	"github.com/mark-rushakoff/go_tftpd/response_agent"
	"github.com/mark-rushakoff/go_tftpd/safe_packets"
)

func TestBegin(t *testing.T) {
	responseAgent := response_agent.MakeMockResponseAgent()

	config := ReadSessionConfig{
		ResponseAgent: responseAgent,
		Reader:        strings.NewReader("Hello!"),
		BlockSize:     512,
	}
	session := NewReadSession(config)
	session.Begin()

	if responseAgent.TotalMessagesSent() != 2 {
		t.Fatalf("Expected 2 messages sent but %v messages were sent", responseAgent.TotalMessagesSent())
	}

	sentAck := responseAgent.MostRecentAck()
	actualBlockNumber := sentAck.BlockNumber
	expectedBlockNumber := uint16(0)
	if actualBlockNumber != expectedBlockNumber {
		t.Errorf("Expected ReadSession to ack with block number %v, received %v", expectedBlockNumber, actualBlockNumber)
		bytes.Equal(nil, nil)
	}

	sentData := responseAgent.MostRecentData()
	actualBlockNumber = sentData.BlockNumber
	expectedBlockNumber = uint16(1)
	if actualBlockNumber != expectedBlockNumber {
		t.Errorf("Expected ReadSession to send data with block number %v, received %v", expectedBlockNumber, actualBlockNumber)
	}

	expectedData := []byte("Hello!")
	if !bytes.Equal(sentData.Data.Data, expectedData) {
		t.Errorf("Expected ReadSession to send data %v, received %v", expectedData, sentData.Data.Data)
	}
}

func TestMultipleDataPackets(t *testing.T) {
	responseAgent := response_agent.MakeMockResponseAgent()

	config := ReadSessionConfig{
		ResponseAgent: responseAgent,
		Reader:        strings.NewReader("12345678abcdef"),
		BlockSize:     8,
	}
	session := NewReadSession(config)
	session.Begin()

	if responseAgent.TotalMessagesSent() != 2 {
		t.Fatalf("Expected 2 messages sent but %v messages were sent", responseAgent.TotalMessagesSent())
	}

	sentData := responseAgent.MostRecentData()
	actualBlockNumber := sentData.BlockNumber
	expectedBlockNumber := uint16(1)
	if actualBlockNumber != expectedBlockNumber {
		t.Errorf("Expected ReadSession to send data with block number %v, received %v", expectedBlockNumber, actualBlockNumber)
	}

	expectedData := []byte("12345678")
	if !bytes.Equal(sentData.Data.Data, expectedData) {
		t.Errorf("Expected ReadSession to send data %v, received %v", expectedData, sentData.Data.Data)
	}

	responseAgent.Reset()
	session.Ack <- safe_packets.NewSafeAck(1)

	sentData = responseAgent.MostRecentData()
	if sentData == nil {
		t.Fatalf("Data not sent")
	}
	actualBlockNumber = sentData.BlockNumber
	expectedBlockNumber = uint16(2)
	if actualBlockNumber != expectedBlockNumber {
		t.Errorf("Expected ReadSession to send data with block number %v, received %v", expectedBlockNumber, actualBlockNumber)
	}

	expectedData = []byte("abcdef")
	if !bytes.Equal(sentData.Data.Data, expectedData) {
		t.Errorf("Expected ReadSession to send data %v, received %v", expectedData, sentData.Data.Data)
	}
}
