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
	bFile, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(bFile, &result); err != nil {
		panic(err)
	}
	return result
}
