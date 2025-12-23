package utils

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"
)

const EnvBeta = "beta"
const EnvProd = "prod"

type HubInfo struct {
	HubId           string `json:"hub_id"`
	SiteId          string `json:"site_id"`
	SiteName        string `json:"site_name"`
	SiteType        string `json:"site_type"`          // mdu,sdu,etc.
	Environment     string `json:"environment"`        // beta / prod
	CloudApiRootUrl string `json:"cloud_api_root_url"` // https://v3.futurehome.io
}

type HubUtils struct {
	hubInfoFilePath string
}

func (cs *HubUtils) SetHubInfoFilePath(hubInfoFilePath string) {
	cs.hubInfoFilePath = hubInfoFilePath
}

func NewHubUtils() *HubUtils {
	return &HubUtils{hubInfoFilePath: "/var/lib/futurehome/hub/hub.json"}
}

func (cs *HubUtils) GetHubInfo() (*HubInfo, error) {
	hubInfo := &HubInfo{}
	configFileBody, err := os.ReadFile(cs.hubInfoFilePath)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(configFileBody, hubInfo)
	if err != nil {
		return nil, err
	}
	return hubInfo, nil
}

func GetNonce() string {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	return fmt.Sprint(r1.Int31())
}
