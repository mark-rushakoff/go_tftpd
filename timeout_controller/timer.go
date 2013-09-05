package timeout_controller

import (
	"time"
)

type timer struct {
	duration time.Duration
	elapsed  chan bool
	restart  chan bool
	destroy  chan bool
}

// 3 calls to reset(duration)
// 1 initial call to stop
// needs to call stop at end

func newTimer(duration time.Duration) *timer {
	t := &timer{
		duration: duration,
		restart:  make(chan bool, 1),
		elapsed:  make(chan bool, 1),
		destroy:  make(chan bool, 1),
	}

	go t.watch()

	return t
}

func (t *timer) Elapsed() <-chan bool {
	return t.elapsed
}

func (t *timer) Reset() {
	t.restart <- true
}

func (t *timer) Destroy() {
	t.destroy <- true
}

func (t *timer) watch() {
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
