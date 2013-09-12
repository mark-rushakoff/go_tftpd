package safepacketprovider

import (
	"testing"
	"time"

	"github.com/mark-rushakoff/go_tftpd/packets"
	"github.com/mark-rushakoff/go_tftpd/test_helpers"
)

var fakeAddr = test_helpers.MakeMockAddr("fake_network", "a")

func TestCanProvideSafeAck(t *testing.T) {
	const blockNum uint16 = 1234
	packetConn := test_helpers.NewMockPacketConnWithBytes(t, fakeAddr, []interface{}{
		uint16(packets.AckOpcode),
		uint16(blockNum),
	})

	provider := NewSafePacketProvider(packetConn)

	go provider.Read()

	select {
	case incomingAck := <-provider.IncomingSafeAck():
		if incomingAck.Ack.BlockNumber != blockNum {
			t.Errorf("Expected ack with block number %v, got %v", blockNum, incomingAck.Ack.BlockNumber)
		}
		if incomingAck.Addr != fakeAddr {
			t.Errorf("Expected ack to have address %v, got %v", fakeAddr, incomingAck.Addr)
		}
	case <-time.After(time.Millisecond):
		t.Fatalf("Did not see SafeAck in time")
	}
}

func TestCanProvideSafeReadRequest(t *testing.T) {
	const blockNum uint16 = 1234
	packetConn := test_helpers.NewMockPacketConnWithBytes(t, fakeAddr, []interface{}{
		uint16(packets.ReadOpcode),
		"foobar",
		byte(0),
		"netascii",
		byte(0),
	})

	provider := NewSafePacketProvider(packetConn)

	go provider.Read()

	select {
	case i := <-provider.IncomingSafeReadRequest():
		if i.Addr != fakeAddr {
			t.Errorf("Expected ack to have address %v, got %v", fakeAddr, i.Addr)
		}
	case <-time.After(time.Millisecond):
		t.Fatalf("Did not see SafeReadRequest in time")
	}
}

func TestCanProvideInvalidMessage(t *testing.T) {
	packetConn := test_helpers.NewMockPacketConnWithBytes(t, fakeAddr, []interface{}{
		uint16(packets.ReadOpcode),
		"foobar",
		byte(0),
		"not a valid mode",
		byte(0),
	})

	provider := NewSafePacketProvider(packetConn)

	go provider.Read()

	select {
	case i := <-provider.IncomingInvalidMessage():
		if i.Addr != fakeAddr {
			t.Errorf("Expected ack to have address %v, got %v", fakeAddr, i.Addr)
		}
		if i.ErrorCode != packets.Undefined {
			t.Errorf("Received error code %v but expected %v", i.ErrorCode, packets.Undefined)
		}
		if i.ErrorMessage != "Invalid mode string" {
			t.Errorf("Received error message '%v' but expected 'Invalid mode string'", i.ErrorMessage)
		}
	case <-time.After(time.Millisecond):
		t.Fatalf("Did not see SafeReadRequest in time")
	}
}
