package timeoutcontroller

import (
	"sync"
	"time"
)

type timer interface {
	Elapsed() <-chan bool
	Restart()
	Destroy()
}

type manualTimer struct {
	duration   time.Duration
	waitFactor uint
	mutex      sync.RWMutex

	elapsed chan bool
	restart chan bool
	destroy chan bool
}

func newTimer(duration time.Duration) timer {
	t := &manualTimer{
		duration:   duration,
		waitFactor: 1,

		restart: make(chan bool, 1),
		elapsed: make(chan bool, 1),
		destroy: make(chan bool, 1),
	}

	go t.watch()

	return t
}

func (t *manualTimer) Elapsed() <-chan bool {
	return t.elapsed
}

func (t *manualTimer) Restart() {
	t.restart <- true

	t.mutex.Lock()
	t.waitFactor = 1
	t.mutex.Unlock()
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
		t.mutex.RLock()
		duration := time.Duration(t.waitFactor) * t.duration
		t.mutex.RUnlock()

		select {
		case <-time.After(duration):
			t.elapsed <- true
			t.mutex.Lock()
			t.waitFactor *= 2
			t.mutex.Unlock()
		case <-t.restart:
			// just restart the loop
		case <-t.destroy:
			return
		}
	}
}
