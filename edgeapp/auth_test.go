package edgeapp

import (
	"testing"
)

func TestFhOAuth2Client_ExchangeRefreshToken(t *testing.T) {

	client := NewFhOAuth2Client("netatmo","auth_test")
	client.SetParameters("tcp://cube.local:1883","","",0,0,0,0)
	err := client.Init()
	if err != nil {
		t.Error("Failed to init sync client")
		t.FailNow()
	}
	t.Log(client.hubToken)
	r , err := client.ExchangeRefreshToken("5ae37d5fac34a5f3d98b7492|ef2b826bc5a80941a5ebfa797762bdc3")
	if err != nil {
		t.Error("Can't fetch token",err)
		t.FailNow()
	}
	t.Log("New access token:"+r.AccessToken)
}