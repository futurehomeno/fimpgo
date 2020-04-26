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
	cf.SaveToFile()

	cf2 := Configs{}
	cf2.path = "./config_test.json"
	cf2.SetCustomConfigs(&AppConf{})
	cf2.LoadFromFile()
	ac2 := cf2.GetCustomConfigs()
	t.Log("Custom conf:",ac2.(*AppConf).Prop1)
}