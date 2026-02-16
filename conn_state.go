package fimpgo

import (
	"fmt"
	"sync"
	"time"
)

type ConnStateT struct {
	mu            sync.Mutex
	connected     chan struct{}
	done          chan struct{}
	onceConnected sync.Once
	onceDone      sync.Once
}

func (c *ConnStateT) Init() {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Always reset for proper restart support
	c.connected = make(chan struct{})
	c.done = make(chan struct{})
	c.onceConnected = sync.Once{}
	c.onceDone = sync.Once{}
}

func (c *ConnStateT) OnConnect() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.onceConnected.Do(func() {
		close(c.connected)
	})
}

func (c *ConnStateT) OnDone() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.done == nil {
		return // not initialized
	}

	c.onceDone.Do(func() {
		close(c.done)
		c.connected = make(chan struct{})
	})
}

func (c *ConnStateT) IsConnected() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	select {
	case <-c.connected:
		return true
	default:
		return false
	}
}

func (c *ConnStateT) WaitConnected(timeout time.Duration) error {
	c.mu.Lock()
	connected := c.connected
	done := c.done
	c.mu.Unlock()

	select {
	case <-time.After(timeout):
		return fmt.Errorf("timeout")
	case <-done:
		return fmt.Errorf("connection_lost")
	case <-connected:
		return nil
	}
}

func (c *ConnStateT) DoneC() <-chan struct{} {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.done
}

func (c *ConnStateT) Lock() {
	c.mu.Lock()
}

func (c *ConnStateT) Unlock() {
	c.mu.Unlock()
}
