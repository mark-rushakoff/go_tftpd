package read_session

import (
	"bytes"
	"strings"
	"testing"

	"github.com/mark-rushakoff/go_tftpd/response_agent"
)

func TestBegin(t *testing.T) {
	responseAgent := response_agent.MakeMockResponseAgent()

	config := ReadSessionConfig{
		ResponseAgent: responseAgent,
		Reader:        strings.NewReader("Hello!"),
		BlockSize:     512,
	}
	session := &ReadSession{config}
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
