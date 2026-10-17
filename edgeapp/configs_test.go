package edgeapp

import (
	"path/filepath"
	"testing"
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
