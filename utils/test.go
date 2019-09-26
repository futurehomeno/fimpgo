package utils

import (
	"encoding/json"
	"io/ioutil"
)

type TestConfig struct {
	BrokerURI string
	SiteId    string

}

func GetTestConfig(path string ) TestConfig {
	var result TestConfig
	bFile, _ := ioutil.ReadFile(path)
	json.Unmarshal(bFile,&result)
	return result
}
