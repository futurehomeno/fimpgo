package edgeapp

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestClient(url string) *FhOAuth2Client {
	c := NewFhOAuth2Client("partner", "app", "")
	c.SetParameters("", url, url, time.Millisecond, 3, 0, 0)
	c.SetHubToken("hub-token")
	return c
}

func TestNewFhOAuth2Client_Env(t *testing.T) {
	prod := NewFhOAuth2Client("p", "a", "prod")
	beta := NewFhOAuth2Client("p", "a", "beta")
	assert.Contains(t, prod.RefreshTokenApiUrl(), "partners.futurehome.io")
	assert.Contains(t, beta.RefreshTokenApiUrl(), "partners-beta.futurehome.io")
	assert.Contains(t, beta.AuthCodeApiUrl(), "partners-beta.futurehome.io")
	assert.Equal(t, "auth_client_a", beta.mqttClientID)
}

func TestFhOAuth2Client_SetParameters(t *testing.T) {
	c := NewFhOAuth2Client("p", "a", "")
	c.SetParameters("tcp://x:1", "auth", "refresh", time.Second, 2, 3, time.Minute)
	assert.Equal(t, "tcp://x:1", c.mqttServerURI)
	assert.Equal(t, "auth", c.AuthCodeApiUrl())
	assert.Equal(t, "refresh", c.RefreshTokenApiUrl())
	assert.Equal(t, time.Second, c.retryDelay)
	assert.Equal(t, 2, c.refreshRetry)
	assert.Equal(t, 3, c.cbRetry)
	assert.Equal(t, time.Minute, c.cbRetryDelay)

	c.SetParameters("", "", "", 0, 0, 0, 0)
	assert.Equal(t, "auth", c.AuthCodeApiUrl())
	assert.Equal(t, 2, c.refreshRetry)

	c.SetAuthCodeApiUrl("a2")
	c.SetRefreshTokenApiUrl("r2")
	assert.Equal(t, "a2", c.AuthCodeApiUrl())
	assert.Equal(t, "r2", c.RefreshTokenApiUrl())
}

func TestFhOAuth2Client_ExchangeRefreshToken_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "Bearer hub-token", r.Header.Get("Authorization"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		var req OAuth2RefreshProxyRequest
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&req))
		assert.Equal(t, OAuth2RefreshProxyRequest{RefreshToken: "rt", PartnerCode: "partner"}, req)
		_, _ = w.Write([]byte(`{"access_token":"at","token_type":"bearer","expires_in":60,"refresh_token":"rt2","scope":"x"}`))
	}))
	defer srv.Close()

	resp, err := newTestClient(srv.URL).ExchangeRefreshToken("rt")
	require.NoError(t, err)
	assert.Equal(t, "at", resp.AccessToken)
	assert.Equal(t, "rt2", resp.RefreshToken)
	assert.Equal(t, int64(60), resp.ExpiresIn)
}

func TestFhOAuth2Client_ExchangeCodeForTokens(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"access_token":"at"}`))
	}))
	defer srv.Close()

	c := newTestClient(srv.URL)
	c.SetAuthCodeApiUrl(srv.URL)
	c.SetRefreshTokenApiUrl(srv.URL)
	resp, err := c.ExchangeCodeForTokens("code")
	require.NoError(t, err)
	assert.Equal(t, "at", resp.AccessToken)
}

func TestFhOAuth2Client_postMsg_Errors(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		hits    int32
	}{
		{"server error retried", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusInternalServerError) }, 3},
		{"invalid json", func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte("nope")) }, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var hits atomic.Int32
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				hits.Add(1)
				tt.handler(w, r)
			}))
			defer srv.Close()

			resp, err := newTestClient(srv.URL).ExchangeRefreshToken("rt")
			require.Error(t, err)
			assert.Nil(t, resp)
			assert.Equal(t, tt.hits, hits.Load())
		})
	}
}

func TestFhOAuth2Client_postMsg_ConnectionError(t *testing.T) {
	srv := httptest.NewServer(http.NotFoundHandler())
	url := srv.URL
	srv.Close()

	_, err := newTestClient(url).ExchangeRefreshToken("rt")
	require.Error(t, err)

	_, err = newTestClient("://bad url").ExchangeRefreshToken("rt")
	require.Error(t, err)
}
