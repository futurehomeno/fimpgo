package edgeapp

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func TestLifecycle_States(t *testing.T) {
	lf := NewAppLifecycle()
	assert.Equal(t, &AppStates{App: AppStateStarting, Connection: ConnStateNA, Config: ConfigStateNotConfigured, Auth: AuthStateNA}, lf.GetAllStates())

	ch := lf.Subscribe("s", 10)
	lf.SetConfigState(ConfigStateConfigured)
	lf.SetAuthState(AuthStateAuthenticated)
	lf.SetConnectionState(ConnStateConnected)
	lf.SetAppState(AppStateRunning, map[string]string{"a": "b"})
	lf.SetLastError("boom")

	assert.Equal(t, State(ConfigStateConfigured), lf.ConfigState())
	assert.Equal(t, State(AuthStateAuthenticated), lf.AuthState())
	assert.Equal(t, State(ConnStateConnected), lf.ConnectionState())
	assert.Equal(t, State(AppStateRunning), lf.AppState())
	assert.Equal(t, State(AppStateStarting), lf.previousAppState)
	assert.Equal(t, "boom", lf.LastError())

	for _, want := range []SystemEvent{
		{Type: SystemEventTypeConfigState, State: ConfigStateConfigured, Info: "sys"},
		{Type: SystemEventTypeAuthState, State: AuthStateAuthenticated, Info: "sys"},
		{Type: SystemEventTypeConnectionState, State: ConnStateConnected, Info: "sys"},
		{Type: SystemEventTypeState, State: AppStateRunning, Info: "sys", Params: map[string]string{"a": "b"}},
	} {
		assert.Equal(t, want, <-ch)
	}
}

func TestLifecycle_PublishAndUnsubscribe(t *testing.T) {
	lf := NewAppLifecycle()
	lf.SetAppState(AppStateRunning, nil)
	ch := lf.Subscribe("s", 1)
	lf.PublishSystemEvent("ev", "src", map[string]string{"k": "v"})
	assert.Equal(t, SystemEvent{Type: SystemEventTypeEvent, Name: "ev", State: AppStateRunning, Params: map[string]string{"k": "v"}}, <-ch)

	lf.Unsubscribe("s")
	lf.PublishSystemEvent("ev2", "src", nil)
	assert.Empty(t, ch)
}

func TestLifecycle_BusyListenerDropsEvents(t *testing.T) {
	lf := NewAppLifecycle()
	ch := lf.Subscribe("busy", 1)
	lf.SetConfigState(ConfigStateInProgress)
	lf.SetAuthState(AuthStateInProgress)
	lf.SetConnectionState(ConnStateConnecting)
	lf.SetAppState(AppStateError, nil)
	lf.PublishSystemEvent("ev", "src", nil)
	require.Len(t, ch, 1)
	assert.Equal(t, SystemEventTypeConfigState, (<-ch).Type)
	assert.Equal(t, State(AppStateError), lf.AppState())
}

func TestLifecycle_WaitForState_Blocking(t *testing.T) {
	lf := NewAppLifecycle()
	lf.SetAppState(AppStateRunning, nil)
	done := make(chan struct{})
	go func() {
		lf.WaitForState("w", SystemEventTypeEvent, AppStateRunning)
		close(done)
	}()
	require.Eventually(t, func() bool {
		lf.busMux.Lock()
		defer lf.busMux.Unlock()
		return len(lf.systemEventBus) == 1
	}, time.Second, time.Millisecond)

	lf.PublishSystemEvent("ev", "src", nil)
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("WaitForState did not return")
	}
	assert.Empty(t, lf.systemEventBus)
}
