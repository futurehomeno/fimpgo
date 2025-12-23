package edgeapp

// TODO: these tests require a running FH broker on cube.local address
/*
func TestFhOAuth2Client_ExchangeRefreshToken(t *testing.T) {
	client := NewFhOAuth2Client("netatmo", "auth_test", "beta")
	client.SetParameters("tcp://cube.local:1883", "", "", 0, 0, 0, 0)
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
}*/
