package fimpgo

import (
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type ConnStateT struct {
	mu            sync.Mutex
	connected     chan struct{}
	done          chan struct{}
	onceConnected sync.Once
	onceDone      sync.Once
}

func (c *ConnStateT) Init() {
	if c.connected != nil {
		log.Warnf("[fimpgo] Already initalized")
		return
	}

	c.mu.Lock()
	c.connected = make(chan struct{})
	c.done = make(chan struct{})
	c.onceConnected = sync.Once{}
	c.onceDone = sync.Once{}
	c.mu.Unlock()
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
	select {
	case <-time.After(timeout):
		return fmt.Errorf("timeout")
	case <-c.connected:
		return nil
	}
}

func (c *ConnStateT) DoneC() <-chan struct{} {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.done
}
