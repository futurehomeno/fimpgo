package edgeapp

import (
	"testing"
	"time"
)

func TestLifecycle_WaitForState(t *testing.T) {
	lf := NewAppLifecycle()
	lf.SetAppState(AppStateRunning,nil)
	lf.WaitForState("test-1",SystemEventTypeState,AppStateRunning)
	lf.SetConfigState(ConfigStateConfigured)
	lf.WaitForState("test-1",SystemEventTypeConfigState,ConfigStateConfigured)
	lf.SetConfigState(ConfigStateConfigured)
	lf.WaitForState("test-1",SystemEventTypeConfigState,ConfigStateConfigured)

	go func() {
		time.Sleep(time.Second*3)
		lf.SetAuthState(AuthStateAuthenticated)
	}()
	lf.WaitForState("test-1",SystemEventTypeAuthState,AuthStateAuthenticated)

	go func() {
		time.Sleep(time.Second*2)
		lf.PublishSystemEvent("testEvent","unitTest",nil)
	}()

	for evt := range lf.Subscribe("test-id",2) {
		if evt.Type == SystemEventTypeEvent && evt.Name == "testEvent"{
			break
		}
	}

	t.Log("OK")

}
