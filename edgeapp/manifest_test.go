package edgeapp

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManifest_SaveLoadLookup(t *testing.T) {
	m := NewManifest()
	m.Configs = []AppConfig{{ID: "c1", ValT: "string", IsRequired: true}}
	m.UIButtons = []UIButton{{ID: "b1"}}
	m.UIBlocks = []AppUBLock{{ID: "u1", Configs: []string{"c1"}}}
	m.InitFlow = []string{"connect"}
	path := filepath.Join(t.TempDir(), "manifest.json")
	require.NoError(t, m.SaveToFile(path))

	loaded := NewManifest()
	require.NoError(t, loaded.LoadFromFile(path))
	assert.Equal(t, m, loaded)

	assert.Equal(t, "c1", loaded.GetAppConfig("c1").ID)
	assert.Equal(t, "b1", loaded.GetButton("b1").ID)
	assert.Equal(t, "u1", loaded.GetUIBlock("u1").ID)
	assert.Nil(t, loaded.GetAppConfig("x"))
	assert.Nil(t, loaded.GetButton("x"))
	assert.Nil(t, loaded.GetUIBlock("x"))
}

func TestManifest_Visibility(t *testing.T) {
	m := &Manifest{Configs: []AppConfig{{ID: "c"}}, UIButtons: []UIButton{{ID: "b"}}, UIBlocks: []AppUBLock{{ID: "u"}}}
	c, b, u := m.GetAppConfig("c"), m.GetButton("b"), m.GetUIBlock("u")
	c.Hide()
	b.Hide()
	u.Hide()
	assert.True(t, m.Configs[0].Hidden && m.UIButtons[0].Hidden && m.UIBlocks[0].Hidden)
	c.Show()
	b.Show()
	u.Show()
	assert.False(t, m.Configs[0].Hidden || m.UIButtons[0].Hidden || m.UIBlocks[0].Hidden)
}

func TestManifest_FileErrors(t *testing.T) {
	dir := t.TempDir()
	m := NewManifest()
	require.Error(t, m.LoadFromFile(filepath.Join(dir, "missing.json")))

	bad := filepath.Join(dir, "bad.json")
	require.NoError(t, os.WriteFile(bad, []byte("{"), 0o600))
	require.Error(t, m.LoadFromFile(bad))

	require.Error(t, m.SaveToFile(filepath.Join(dir, "nodir", "m.json")))

	m.ConfigState = make(chan int)
	require.Error(t, m.SaveToFile(filepath.Join(dir, "m.json")))
}
