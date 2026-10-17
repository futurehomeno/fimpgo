package edgeapp

import (
	"strconv"
	"testing"
	"time"
)

func TestLifecycle_SubscribeWhileStatesChange(t *testing.T) {
	lf := NewAppLifecycle()
	done := make(chan struct{})

	go func() {
		defer close(done)

		for i := range 10000 {
			id := strconv.Itoa(i % 10)
			lf.Subscribe(id, 1)
			lf.Unsubscribe(id)
		}
	}()

	for {
		select {
		case <-done:
			return
		default:
			lf.SetConfigState(ConfigStateConfigured)
			lf.SetAuthState(AuthStateAuthenticated)
			lf.SetConnectionState(ConnStateConnected)
		}
	}
}

func TestLifecycle_WaitForStateChangedWhileStarting(t *testing.T) {
	lf := NewAppLifecycle()
	returned := make(chan struct{})

	lf.busMux.Lock()

	go func() {
		lf.WaitForState("test-1", SystemEventTypeState, AppStateRunning)
		close(returned)
	}()

	// Lets the wait read the old state, then changes it without an event the wait could receive.
	time.Sleep(100 * time.Millisecond)
	lf.appState = AppStateRunning
	lf.busMux.Unlock()

	select {
	case <-returned:
	case <-time.After(time.Second):
		t.Fatal("wait missed the state change")
	}
}

func TestLifecycle_WaitForState(t *testing.T) {
	lf := NewAppLifecycle()
	lf.SetAppState(AppStateRunning, nil)
	lf.WaitForState("test-1", SystemEventTypeState, AppStateRunning)
	lf.SetConfigState(ConfigStateConfigured)
	lf.WaitForState("test-1", SystemEventTypeConfigState, ConfigStateConfigured)
	lf.SetConfigState(ConfigStateConfigured)
	lf.WaitForState("test-1", SystemEventTypeConfigState, ConfigStateConfigured)

	go func() {
		time.Sleep(time.Second * 3)
		lf.SetAuthState(AuthStateAuthenticated)
	}()
	lf.WaitForState("test-1", SystemEventTypeAuthState, AuthStateAuthenticated)

	go func() {
		time.Sleep(time.Second * 2)
		lf.PublishSystemEvent("testEvent", "unitTest", nil)
	}()

	for evt := range lf.Subscribe("test-id", 2) {
		if evt.Type == SystemEventTypeEvent && evt.Name == "testEvent" {
			break
		}
	}
}
