package testhelpers

type MockAddr struct {
	network string
	str     string
}

func MakeMockAddr(network string, str string) *MockAddr {
	return &MockAddr{
		network: network,
		str:     str,
	}
}

func (a *MockAddr) Network() string {
	return a.network
}

func (a *MockAddr) String() string {
	return a.str
}
