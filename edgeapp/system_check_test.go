package edgeapp

import (
	"testing"
	"time"
)
import log "github.com/sirupsen/logrus"

func TestSystemCheck_IsNetworkAvailable(t *testing.T) {
	log.SetLevel(log.TraceLevel)
	sc := NewSystemCheck()
	if !sc.IsNetworkAvailable() {
		t.Error("Network is not available")
	}
}

func TestSystemCheck_IsInternetAvailable(t *testing.T) {
	log.SetLevel(log.TraceLevel)
	sc := NewSystemCheck()
	if !sc.IsInternetAvailable() {
		t.Error("Internet is not available")
	}
}

func TestSystemCheck_WaitForInternet(t *testing.T) {
	log.SetLevel(log.TraceLevel)
	sc := NewSystemCheck()
	if err := sc.WaitForInternet(15*time.Second);err != nil {
		t.Error("Internet is not available")
	}
}