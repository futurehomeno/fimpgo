package utils

import (
	"net"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRouteIncludesTopic(t *testing.T) {
	for _, c := range []struct {
		route, topic string
		want         bool
	}{
		{"a/b/c", "a/b/c", true},
		{"a/+/c", "a/x/c", true},
		{"a/+/c", "a/x/d", false},
		{"a/#", "a/b/c", true},
		{"a/#", "a", true},
		{"#", "x/y", true},
		{"a/b", "a", false},
		{"a", "a/b", false},
	} {
		assert.Equal(t, c.want, RouteIncludesTopic(c.route, c.topic), "%s vs %s", c.route, c.topic)
	}
}

func TestFileExistsAndCopyFile(t *testing.T) {
	dir := t.TempDir()
	src, dst := filepath.Join(dir, "src"), filepath.Join(dir, "dst")
	require.NoError(t, os.WriteFile(src, []byte("data"), 0600))

	assert.True(t, FileExists(src))
	assert.False(t, FileExists(dir))
	assert.False(t, FileExists(dst))

	require.NoError(t, CopyFile(src, dst))
	got, err := os.ReadFile(dst) //nolint:gosec
	require.NoError(t, err)
	assert.Equal(t, "data", string(got))

	assert.Error(t, CopyFile(filepath.Join(dir, "none"), dst))
	assert.Error(t, CopyFile(dir, dst))
	assert.Error(t, CopyFile(src, filepath.Join(dir, "no", "dst")))
}

func TestGetNonce(t *testing.T) {
	assert.Regexp(t, `^\d+$`, GetNonce())
}

func TestHubUtils_GetHubInfo(t *testing.T) {
	path := filepath.Join(t.TempDir(), "hub.json")
	h := NewHubUtils()
	assert.Equal(t, "/var/lib/futurehome/hub/hub.json", h.hubInfoFilePath)
	h.SetHubInfoFilePath(path)

	_, err := h.GetHubInfo()
	assert.Error(t, err)

	require.NoError(t, os.WriteFile(path, []byte("{"), 0600))
	_, err = h.GetHubInfo()
	assert.Error(t, err)

	require.NoError(t, os.WriteFile(path, []byte(`{"hub_id":"h1","site_id":"s1","environment":"beta"}`), 0600))
	info, err := h.GetHubInfo()
	require.NoError(t, err)
	assert.Equal(t, &HubInfo{HubId: "h1", SiteId: "s1", Environment: EnvBeta}, info)
}

func TestRoutedInterface(t *testing.T) {
	_, err := RoutedInterface("tcp", 0)
	assert.Error(t, err)

	ifi, err := RoutedInterface("ip4", net.FlagLoopback)
	if err == nil {
		assert.NotEmpty(t, ifi.Name)
	}
	_, err = RoutedInterface("ip4", net.FlagLoopback|net.FlagPointToPoint|net.FlagBroadcast)
	assert.Error(t, err)
}

func TestRoutableIP(t *testing.T) {
	for _, c := range []struct {
		network, ip string
		want        bool
	}{
		{"ip4", "127.0.0.1", true},
		{"ip4", "10.0.0.1", true},
		{"ip4", "fe80::1", false},
		{"ip6", "fe80::1", true},
		{"ip6", "::1", false},
		{"ip6", "10.0.0.1", false},
		{"ip", "::1", false},
		{"ip", "10.0.0.1", true},
		{"ip", "2001:db8::1", true},
		{"ip4", "0.0.0.0", false},
		{"ip4", "224.0.0.1", false},
	} {
		_, ok := routableIP(c.network, net.ParseIP(c.ip))
		assert.Equal(t, c.want, ok, "%s %s", c.network, c.ip)
	}
}

func TestHasRoutableIP(t *testing.T) {
	lo, err := net.InterfaceByName("lo")
	if err != nil {
		t.Skip("no lo interface")
	}
	ip, ok := hasRoutableIP("ip4", lo)
	assert.True(t, ok)
	assert.True(t, ip.IsLoopback())
	_, ok = hasRoutableIP("ip", lo)
	assert.False(t, ok)
}
