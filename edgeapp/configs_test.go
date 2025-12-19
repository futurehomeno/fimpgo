package edgeapp

import "testing"

type AppConf struct {
	Prop1 string
}

func TestConfigs_LoadFromFile(t *testing.T) {
	cf := Configs{}
	cf.path = "./config_test.json"
	ac := AppConf{Prop1: "39"}
	cf.SetCustomConfigs(ac)
	if err := cf.SaveToFile(); err != nil {
		t.Error("Failed to save config:", err)
		return
	}

	cf2 := Configs{}
	cf2.path = "./config_test.json"
	cf2.SetCustomConfigs(&AppConf{})
	if err := cf2.LoadFromFile(); err != nil {
		t.Error("Failed to load config:", err)
		return
	}
	ac2 := cf2.GetCustomConfigs()
	t.Log("Custom conf:", ac2.(*AppConf).Prop1)
}
