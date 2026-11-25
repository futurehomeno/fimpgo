package edgeapp

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/futurehomeno/fimpgo/utils"
)

// TODO: these tests require a running FH broker on cube.local address
func TestFhOAuth2Client_ExchangeRefreshToken(t *testing.T) {
	t.Skip()
	client := NewFhOAuth2Client("netatmo", "auth_test", "beta")
	client.SetParameters("tcp://cube.local:11883", "", "", 0, 0, 0, 0)
	err := client.Init()
	if err != nil {
		t.Error("Failed to init sync client")
		t.FailNow()
	}
	_, err = client.ExchangeRefreshToken("the token must be set here")
	if err != nil {
		t.Error("Can't fetch token", err)
		t.FailNow()
	}
}

func TestFhOAuth2Client_LoadHubTokenFromCBAgain(t *testing.T) {
	client := NewFhOAuth2Client("p", "auth_reload_test", "")
	client.SetParameters("tcp://127.0.0.1:11883", "", "", 0, 0, 1, time.Second)

	assert.ErrorIs(t, client.LoadHubTokenFromCB(), utils.ErrTimeout)
	assert.ErrorIs(t, client.LoadHubTokenFromCB(), utils.ErrTimeout)
}

func TestFhOAuth2Client_LoadHubTokenFromCBReleasesOnSubscribeError(t *testing.T) {
	client := NewFhOAuth2Client("p", "x/#", "")
	client.SetParameters("tcp://127.0.0.1:11883", "", "", 0, 0, 1, time.Second)

	assert.Error(t, client.LoadHubTokenFromCB())
	assert.Nil(t, client.mqt)
	assert.Nil(t, client.syncClient)
}

func TestFhOAuth2Client_LoadHubTokenFromCBConcurrent(t *testing.T) {
	client := NewFhOAuth2Client("p", "auth_concurrent_test", "")
	client.SetParameters("tcp://127.0.0.1:11883", "", "", 0, 0, 1, time.Second)

	errs := make(chan error, 2)
	for range 2 {
		go func() { errs <- client.LoadHubTokenFromCB() }()
	}

	assert.ErrorIs(t, <-errs, utils.ErrTimeout)
	assert.ErrorIs(t, <-errs, utils.ErrTimeout)
}
