package timeoutcontroller

import (
	"time"
)

type timer interface {
	Elapsed() <-chan bool
	Restart()
	Destroy()
}

type manualTimer struct {
	duration time.Duration
	elapsed  chan bool
	restart  chan bool
	destroy  chan bool
}

func newTimer(duration time.Duration) timer {
	t := &manualTimer{
		duration: duration,
		restart:  make(chan bool, 1),
		elapsed:  make(chan bool, 1),
		destroy:  make(chan bool, 1),
	}

	go t.watch()

	return t
}

func (t *manualTimer) Elapsed() <-chan bool {
	return t.elapsed
}

func (t *manualTimer) Restart() {
	t.restart <- true
}

func (t *manualTimer) Destroy() {
	t.destroy <- true
}

func (t *manualTimer) watch() {
	select {
	case <-t.restart:
		// need initial restart call to get going
	}

	for {
		select {
		case <-time.After(t.duration):
			t.elapsed <- true
		case <-t.restart:
			// just restart the loop
		case <-t.destroy:
			return
		}
	}
}
