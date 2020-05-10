package edgeapp

import (
	"testing"
)

func TestFhOAuth2Client_ExchangeRefreshToken(t *testing.T) {
	client := &FhOAuth2Client{partnerName: "netatmo", mqttServerURI: "tcp://cube.local:1883",
		mqttClientID: "fhouthclient",
		refreshTokenApiUrl: "https://partners-beta.futurehome.io/api/control/edge/proxy/refresh"}
	client.retryDelay = 10
	client.refreshRetry = 3
	client.cbRetryDelay = 30
	client.cbRetry = 7
	err := client.ConfigureFimpSyncClient()
	if err != nil {
		t.Error("Failed to init sync client")
		t.FailNow()
	}
	err =  client.LoadHubTokenFromCB()
	if err != nil {
		t.Error("Failed to load token from CB")
		t.FailNow()
	}
	t.Log(client.hubToken)
	r , err := client.ExchangeRefreshToken("")
	if err != nil {
		t.Error("Can't fetch token",err)
		t.FailNow()
	}
	t.Log("New access token:"+r.AccessToken)
}