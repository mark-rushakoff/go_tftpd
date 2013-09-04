package timeout_controller

import (
	"sync"
)

type MockTimeoutController struct {
	countdownCalls uint
	restartCalls   uint
	timeout        chan bool

	mutex sync.RWMutex
}

func MakeMockTimeoutController() *MockTimeoutController {
	return &MockTimeoutController{
		timeout: make(chan bool, 1),
	}
}

func (c *MockTimeoutController) Timeout() chan bool {
	return c.timeout
}

func (c *MockTimeoutController) CauseExpiredTimeout() {
	c.timeout <- true
}

func (c *MockTimeoutController) CauseNonExpiredTimeout() {
	c.timeout <- false
}

func (c *MockTimeoutController) Countdown() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.countdownCalls++
	return nil
}

func (c *MockTimeoutController) CountdownCalls() uint {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.countdownCalls
}

func (c *MockTimeoutController) Restart() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.restartCalls++
}

func (c *MockTimeoutController) RestartCalls() uint {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.restartCalls
}

func (c *MockTimeoutController) Stop() {
}
