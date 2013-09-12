package timeoutcontroller

type mockTimer struct {
	elapsed chan bool
	restart chan<- bool
	destroy chan<- bool
}

func NewMockTimer(restart chan<- bool, destroy chan<- bool) *mockTimer {
	return &mockTimer{
		restart: restart,
		destroy: destroy,
		elapsed: make(chan bool),
	}
}

func (t *mockTimer) Elapsed() <-chan bool {
	return t.elapsed
}

func (t *mockTimer) Elapse() {
	t.elapsed <- true
}

func (t *mockTimer) Restart() {
	t.restart <- true
}

func (t *mockTimer) Destroy() {
	t.destroy <- true
}
