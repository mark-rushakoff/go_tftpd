package readsession

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/mark-rushakoff/go_tftpd/safepackets"
)

func TestBegin(t *testing.T) {
	dataChan := make(chan *safepackets.SafeData, 1)
	handler := &PluggableHandler{
		SendDataHandler: func(d *safepackets.SafeData) {
			dataChan <- d
		},
	}
	config := &Config{
		Reader:    strings.NewReader("foobar"),
		BlockSize: 2,
	}
	session := NewReadSession(config, handler, func() {})
	session.Begin()

	select {
	case d := <-dataChan:
		if d.BlockNumber != 1 {
			t.Errorf("Expected block number 1, got %v", d.BlockNumber)
		}
		if !bytes.Equal(d.Data.Data, []byte("fo")) {
			t.Errorf("Expected data packet of 'fo', saw %v", d.Data)
		}
	case <-time.After(time.Millisecond):
		t.Fatalf("Did not see data packet in time")
	}
}

func TestCurrentAckAdvancesData(t *testing.T) {
	dataChan := make(chan *safepackets.SafeData, 1)
	handler := &PluggableHandler{
		SendDataHandler: func(d *safepackets.SafeData) {
			dataChan <- d
		},
	}
	config := &Config{
		Reader:    strings.NewReader("foobar"),
		BlockSize: 2,
	}
	session := NewReadSession(config, handler, func() {})
	session.Begin()

	select {
	case d := <-dataChan:
		// same case asserted in another test
		if d.BlockNumber != 1 {
			t.Errorf("Expected block number 1, got %v", d.BlockNumber)
		}
	case <-time.After(time.Millisecond):
		t.Fatalf("Did not see data packet in time")
	}

	session.HandleAck(safepackets.NewSafeAck(1))

	select {
	case d := <-dataChan:
		// same case asserted in another test
		if d.BlockNumber != 2 {
			t.Errorf("Expected block number 2, got %v", d.BlockNumber)
		}
		if !bytes.Equal(d.Data.Data, []byte("ob")) {
			t.Errorf("Expected data packet of 'ob', saw %v", d.Data)
		}
	case <-time.After(time.Millisecond):
		t.Fatalf("Did not see data packet in time")
	}
}

func TestPreviousAckRepeatsData(t *testing.T) {
	dataChan := make(chan *safepackets.SafeData, 1)
	handler := &PluggableHandler{
		SendDataHandler: func(d *safepackets.SafeData) {
			dataChan <- d
		},
	}
	config := &Config{
		Reader:    strings.NewReader("foobar"),
		BlockSize: 2,
	}
	session := NewReadSession(config, handler, func() {})

	session.Begin()
	select {
	case d := <-dataChan:
		// same case asserted in another test
		if d.BlockNumber != 1 {
			t.Errorf("Expected block number 1, got %v", d.BlockNumber)
		}
	case <-time.After(time.Millisecond):
		t.Fatalf("Did not see data packet in time")
	}

	session.HandleAck(safepackets.NewSafeAck(1))
	select {
	case d := <-dataChan:
		// same case asserted in another test
		if d.BlockNumber != 2 {
			t.Errorf("Expected block number 2, got %v", d.BlockNumber)
		}
	case <-time.After(time.Millisecond):
		t.Fatalf("Did not see data packet in time")
	}

	session.HandleAck(safepackets.NewSafeAck(1))
	select {
	case d := <-dataChan:
		// same case asserted in another test
		if d.BlockNumber != 2 {
			t.Errorf("Expected block number 2, got %v", d.BlockNumber)
		}
		if !bytes.Equal(d.Data.Data, []byte("ob")) {
			t.Errorf("Expected data packet of 'ob', saw %v", d.Data)
		}
	case <-time.After(time.Millisecond):
		t.Fatalf("Did not see data packet in time")
	}
}

func TestResendRepeatsData(t *testing.T) {
	dataChan := make(chan *safepackets.SafeData, 1)
	handler := &PluggableHandler{
		SendDataHandler: func(d *safepackets.SafeData) {
			dataChan <- d
		},
	}
	config := &Config{
		Reader:    strings.NewReader("foobar"),
		BlockSize: 2,
	}
	session := NewReadSession(config, handler, func() {})

	session.Begin()
	select {
	case d := <-dataChan:
		// same case asserted in another test
		if d.BlockNumber != 1 {
			t.Errorf("Expected block number 1, got %v", d.BlockNumber)
		}
	case <-time.After(time.Millisecond):
		t.Fatalf("Did not see data packet in time")
	}

	session.Resend()
	select {
	case d := <-dataChan:
		// same case asserted in another test
		if d.BlockNumber != 1 {
			t.Errorf("Expected block number 1, got %v", d.BlockNumber)
		}
	default:
		t.Fatalf("Did not see data packet in time")
	}

	session.HandleAck(safepackets.NewSafeAck(1))
	select {
	case d := <-dataChan:
		if d.BlockNumber != 2 {
			t.Errorf("Expected block number 2, got %v", d.BlockNumber)
		}
		if !bytes.Equal(d.Data.Data, []byte("ob")) {
			t.Errorf("Expected ob, saw %v", d.Data.Data)
		}
	default:
		t.Fatalf("Did not see data packet in time")
	}

	session.Resend()
	select {
	case d := <-dataChan:
		if d.BlockNumber != 2 {
			t.Errorf("Expected block number 2, got %v", d.BlockNumber)
		}
		if !bytes.Equal(d.Data.Data, []byte("ob")) {
			t.Errorf("Expected ob, saw %v", d.Data.Data)
		}
	default:
		t.Fatalf("Did not see data packet in time")
	}
}

func TestFinishesInSinglePacket(t *testing.T) {
	dataChan := make(chan *safepackets.SafeData, 1)
	finished := make(chan bool, 1)
	handler := &PluggableHandler{
		SendDataHandler: func(d *safepackets.SafeData) {
			dataChan <- d
		},
	}
	config := &Config{
		Reader:    strings.NewReader("foobar"),
		BlockSize: 24,
	}
	session := NewReadSession(config, handler, func() {
		finished <- true
	})
	session.Begin()

	select {
	case d := <-dataChan:
		if d.BlockNumber != 1 {
			// same case asserted in another test
			t.Errorf("Expected block number 1, got %v", d.BlockNumber)
		}
	case <-time.After(time.Millisecond):
		t.Fatalf("Did not see data packet in time")
	}

	select {
	case <-finished:
		t.Errorf("Expected session not to be finished before last ack arrived")
	default:
		// ok
	}

	session.HandleAck(safepackets.NewSafeAck(1))

	select {
	case <-finished:
		// ok
	default:
		t.Errorf("Expected session to be finished after last ack arrived")
	}
}

func TestFinishesInMultiplePackets(t *testing.T) {
	dataChan := make(chan *safepackets.SafeData, 1)
	finished := make(chan bool, 1)
	handler := &PluggableHandler{
		SendDataHandler: func(d *safepackets.SafeData) {
			dataChan <- d
		},
	}
	config := &Config{
		Reader:    strings.NewReader("foobar"),
		BlockSize: 5,
	}
	session := NewReadSession(config, handler, func() {
		finished <- true
	})
	session.Begin()

	select {
	case d := <-dataChan:
		if d.BlockNumber != 1 {
			t.Errorf("Expected block number 1, got %v", d.BlockNumber)
		}

		if !bytes.Equal(d.Data.Data, []byte("fooba")) {
			t.Errorf("Expected fooba, saw %v", d.Data.Data)
		}
	case <-time.After(time.Millisecond):
		t.Fatalf("Did not see data packet in time")
	}

	select {
	case <-finished:
		t.Errorf("Expected session not to be finished before last ack arrived")
	default:
		// ok
	}

	session.HandleAck(safepackets.NewSafeAck(1))
	select {
	case d := <-dataChan:
		if d.BlockNumber != 2 {
			t.Errorf("Expected block number 2, got %v", d.BlockNumber)
		}

		if !bytes.Equal(d.Data.Data, []byte("r")) {
			t.Errorf("Expected r, saw %v", d.Data.Data)
		}
	case <-time.After(time.Millisecond):
		t.Fatalf("Did not see data packet in time")
	}

	session.HandleAck(safepackets.NewSafeAck(2))

	select {
	case <-finished:
		// ok
	default:
		t.Errorf("Expected session to be finished after last ack arrived")
	}
}
