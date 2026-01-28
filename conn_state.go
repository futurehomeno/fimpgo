package fimpgo

import (
	"fmt"
	"sync"
	"time"
)

type ConnStateT struct {
	mu          sync.Mutex
	started     chan struct{}
	done        chan struct{}
	onceStarted sync.Once
	onceDone    sync.Once
}

func (c *ConnStateT) Init() {
	c.mu.Lock()
	c.started = make(chan struct{})
	c.done = make(chan struct{})
	c.onceStarted = sync.Once{}
	c.onceDone = sync.Once{}
	c.mu.Unlock()
}

func (c *ConnStateT) OnConnect() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.onceStarted.Do(func() {
		close(c.started)
	})
}

func (c *ConnStateT) OnDone() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.onceDone.Do(func() {
		close(c.done)
	})
}

func (c *ConnStateT) IsConnected() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	select {
	case <-c.started:
		return true
	default:
		return false
	}
}

func (c *ConnStateT) WaitConnected(timeout time.Duration) error {
	c.mu.Lock()
	ch := c.started
	c.mu.Unlock()

	select {
	case <-time.After(timeout):
		return fmt.Errorf("timeout")
	case <-ch:
		return nil
	}
}

func (c *ConnStateT) DoneC() <-chan struct{} {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.done
}
