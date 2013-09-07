package timeout_controller

import (
	"sync"
)

type tryCounter struct {
	triesRemaining uint
	tryLimit       uint
	lock           sync.RWMutex
}

func (c *tryCounter) Decrement() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.triesRemaining--
}

func (c *tryCounter) IsZero() bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.triesRemaining == 0
}

func (c *tryCounter) Reset() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.triesRemaining = c.tryLimit
}
