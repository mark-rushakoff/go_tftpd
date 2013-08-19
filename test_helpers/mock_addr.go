package test_helpers

type MockAddr struct {
}

func (a *MockAddr) Network() string {
	return ""
}

func (a *MockAddr) String() string {
	return ""
}
