package safepackets

import (
	"testing"

	"github.com/mark-rushakoff/go_tftpd/packets"
)

func TestSafeAckConversion(t *testing.T) {
	ack := &packets.Ack{BlockNumber: 16}
	safeAck := NewConverter().FromAck(ack)
	if safeAck.BlockNumber != 16 {
		t.Fatalf("FromAck converted block number incorrectly")
	}
}

func TestSafeReadConversion(t *testing.T) {
	type testCase struct {
		actualFilename string
		actualMode     string

		expectedMode ReadWriteMode
	}

	testCases := []testCase{
		{"foo", "netascii", NetAscii},
		{"foo", "NetAscii", NetAscii},
		{"foo", "octet", Octet},
		{"foo", "Octet", Octet},
	}

	for _, testCase := range testCases {
		read := &packets.ReadRequest{
			Filename: testCase.actualFilename,
			Mode:     testCase.actualMode,
			Options:  nil,
		}

		safeRead, err := NewConverter().FromReadRequest(read)

		if err != nil {
			t.Fatalf("ReadRequest should not have caused error in conversion")
		}

		if safeRead.Filename != testCase.actualFilename {
			t.Fatalf("Did not convert filename correctly")
		}

		if safeRead.Mode != testCase.expectedMode {
			t.Fatalf("Did not convert mode correctly")
		}
	}
}

func TestSafeReadConversionCanReturnError(t *testing.T) {
	read := &packets.ReadRequest{
		Filename: "foo",
		Mode:     "baz",
		Options:  nil,
	}

	_, err := NewConverter().FromReadRequest(read)

	if err == nil {
		t.Fatalf("ReadRequest should have caused error in conversion")
	}

	if err.Code() != packets.Undefined {
		t.Errorf("Incorrect error code from invalid mode")
	}
}
