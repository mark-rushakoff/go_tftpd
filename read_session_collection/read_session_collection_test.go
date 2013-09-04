package session_manager

import (
	"testing"

	"github.com/mark-rushakoff/go_tftpd/read_session"
	"github.com/mark-rushakoff/go_tftpd/test_helpers"
)

func TestAddSessionCallsBegin(t *testing.T) {
	begin := make(chan bool, 1)
	session := &read_session.MockReadSession{
		BeginHandler: func() {
			begin <- true
		},
	}
	addr := test_helpers.MakeMockAddr("fake_net", "a")

	manager := NewReadSessionCollection()

	manager.Add(session, addr)

	select {
	case <-begin:
	// ok
	default:
		t.Errorf("Add did not call session.Begin")
	}
}

func TestAddSessionMakesFetchable(t *testing.T) {
	session := &read_session.MockReadSession{
		BeginHandler: func() {
		},
	}
	addr := test_helpers.MakeMockAddr("fake_network", "a")

	manager := NewReadSessionCollection()
	manager.Add(session, addr)

	s, ok := manager.Fetch(addr)
	if !ok {
		t.Fatalf("Should have been able to fetch session")
	}
	if s != session {
		t.Fatalf("Incorrect session returned")
	}
}

func TestRemoveMakesFetchFail(t *testing.T) {
	session := &read_session.MockReadSession{
		BeginHandler: func() {
		},
	}
	addr := test_helpers.MakeMockAddr("fake_network", "a")

	manager := NewReadSessionCollection()
	manager.Add(session, addr)
	manager.Remove(addr)

	_, ok := manager.Fetch(addr)
	if ok {
		t.Fatalf("Should not have been able to fetch removed session")
	}
}
