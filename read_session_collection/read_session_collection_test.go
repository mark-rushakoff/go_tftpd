package read_session_collection

import (
	"testing"

	"github.com/mark-rushakoff/go_tftpd/test_helpers"
	"github.com/mark-rushakoff/go_tftpd/timeout_controller"
)

func TestAddSessionMakesFetchable(t *testing.T) {
	session := &timeout_controller.MockTimeoutController{}
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
	session := &timeout_controller.MockTimeoutController{}
	addr := test_helpers.MakeMockAddr("fake_network", "a")

	manager := NewReadSessionCollection()
	manager.Add(session, addr)
	manager.Remove(addr)

	_, ok := manager.Fetch(addr)
	if ok {
		t.Fatalf("Should not have been able to fetch removed session")
	}
}