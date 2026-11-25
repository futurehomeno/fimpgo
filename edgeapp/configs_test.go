package edgeapp

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type AppConf struct {
	Prop1 string
}

func TestConfigs_LoadFromFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	cf := Configs{}
	cf.path = path
	ac := AppConf{Prop1: "39"}
	cf.SetCustomConfigs(ac)
	if err := cf.SaveToFile(); err != nil {
		t.Error("Failed to save config:", err)
		return
	}

	cf2 := Configs{}
	cf2.path = path
	cf2.SetCustomConfigs(&AppConf{})
	if err := cf2.LoadFromFile(); err != nil {
		t.Error("Failed to load config:", err)
		return
	}
}

func newWorkDir(t *testing.T, defaults string) string {
	t.Helper()
	dir := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "data"), 0o750))
	if defaults != "" {
		require.NoError(t, os.MkdirAll(filepath.Join(dir, "defaults"), 0o750))
		require.NoError(t, os.WriteFile(filepath.Join(dir, "defaults", "config.json"), []byte(defaults), 0o600))
	}
	return dir
}

func TestNewConfigs(t *testing.T) {
	dir := newWorkDir(t, `{"mqtt_server_uri":"tcp://default:1883","log_level":"info"}`)
	cf, err := NewConfigs(dir)
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(dir, "data"), cf.GetDataDir())
	assert.Equal(t, filepath.Join(dir, "defaults"), cf.GetDefaultDir())

	require.NoError(t, cf.LoadFromFile())
	assert.Equal(t, "tcp://default:1883", cf.MqttServerURI)

	cf.MqttServerURI = "tcp://changed:1883"
	cf.SetCustomConfigs(map[string]any{"k": "v"})
	require.NoError(t, cf.SaveToFile())
	assert.Equal(t, "auto", cf.ConfiguredBy)
	assert.NotEmpty(t, cf.ConfiguredAt)

	cf2, err := NewConfigs(dir)
	require.NoError(t, err)
	require.NoError(t, cf2.LoadFromFile())
	assert.Equal(t, "tcp://changed:1883", cf2.MqttServerURI)
	assert.Equal(t, map[string]any{"k": "v"}, cf2.GetCustomConfigs())

	require.NoError(t, cf2.LoadDefaults())
	require.NoError(t, cf2.LoadFromFile())
	assert.Equal(t, "tcp://default:1883", cf2.MqttServerURI)
}

func TestConfigs_Errors(t *testing.T) {
	_, err := NewConfigs(newWorkDir(t, ""))
	require.Error(t, err)

	dir := newWorkDir(t, "{")
	cf, err := NewConfigs(dir)
	require.NoError(t, err)
	require.Error(t, cf.LoadFromFile())

	require.NoError(t, os.Remove(cf.path))
	require.Error(t, cf.LoadFromFile())
	require.NoError(t, cf.LoadDefaults())

	cf.path = filepath.Join(dir, "nodir", "config.json")
	require.Error(t, cf.SaveToFile())

	cf.path = filepath.Join(dir, "data", "config.json")
	cf.CustomConfigs = make(chan int)
	require.Error(t, cf.SaveToFile())

	require.NoError(t, os.Remove(filepath.Join(dir, "defaults", "config.json")))
	require.Error(t, cf.LoadDefaults())
}
