package utils

import (
	"encoding/json"
	"os"
)

type TestConfig struct {
	BrokerURI string
	SiteId    string
}

func GetTestConfig(path string) TestConfig {
	var result TestConfig
	bFile, _ := os.ReadFile(path)
	json.Unmarshal(bFile, &result)
	return result
}
