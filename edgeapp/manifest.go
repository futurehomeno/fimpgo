package edgeapp

import (
	"encoding/json"
	"os"

	"github.com/futurehomeno/fimpgo/fimptype"
)

type Manifest struct {
	Configs     []AppConfig  `json:"configs"`
	UIBlocks    []AppUBLock  `json:"ui_blocks"`
	UIButtons   []UIButton   `json:"ui_buttons"`
	Auth        AppAuth      `json:"auth"`
	InitFlow    []string     `json:"init_flow"`
	Services    []AppService `json:"services"`
	AppState    AppStates    `json:"app_state"`
	ConfigState any          `json:"config_state"`
}

type AppConfig struct {
	ID          string            `json:"id"`
	Label       MultilingualLabel `json:"label"`
	ValT        string            `json:"val_t"`
	UI          AppConfigUI       `json:"ui"`
	Val         Value             `json:"val"`
	IsRequired  bool              `json:"is_required"`
	ConfigPoint string            `json:"config_point"`
	Hidden      bool              `json:"hidden"` //
}

func (b *AppConfig) Hide() {
	b.Hidden = true
}

func (b *AppConfig) Show() {
	b.Hidden = true
}

type MultilingualLabel map[string]string

type AppAuth struct {
	Type                  string `json:"type"`
	CodeGrantLoginPageUrl string `json:"code_grant_login_page_url"`
	RedirectURL           string `json:"redirect_url"`
	ClientID              string `json:"client_id"`
	Secret                string `json:"secret"`
	PartnerID             string `json:"partner_id"`
	AuthEndpoint          string `json:"auth_endpoint"`
}

type AppService struct {
	Name       string               `json:"name"`
	Alias      string               `json:"alias"`
	Address    string               `json:"address"`
	Interfaces []fimptype.Interface `json:"interfaces"`
}

type Value struct {
	Default any `json:"default"`
}

type AppConfigUI struct {
	Type   string `json:"type"`
	Select any    `json:"select"`
}

type UIButton struct {
	ID    string            `json:"id"`
	Label MultilingualLabel `json:"label"`
	Req   struct {
		Serv  string `json:"serv"`
		IntfT string `json:"intf_t"`
		Val   string `json:"val"`
	} `json:"req"`
	Hidden bool `json:"hidden"`
}

func (b *UIButton) Hide() {
	b.Hidden = true
}

func (b *UIButton) Show() {
	b.Hidden = true
}

type ButtonActionResponse struct {
	Operation       string `json:"op"`
	OperationStatus string `json:"op_status"`
	Next            string `json:"next"`
	ErrorCode       string `json:"error_code"`
	ErrorText       string `json:"error_text"`
}

type AppUBLock struct {
	ID      string            `json:"id"`
	Header  MultilingualLabel `json:"header"`
	Text    MultilingualLabel `json:"text"`
	Configs []string          `json:"configs"`
	Buttons []string          `json:"buttons"`
	Footer  MultilingualLabel `json:"footer"`
	Hidden  bool              `json:"hidden"`
}

func (b *AppUBLock) Hide() {
	b.Hidden = true
}

func (b *AppUBLock) Show() {
	b.Hidden = true
}

func NewManifest() *Manifest {
	return &Manifest{}
}

func (m *Manifest) LoadFromFile(filePath string) error {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(file, m)
	if err != nil {
		return err
	}
	return nil
}

func (m *Manifest) SaveToFile(filePath string) error {
	flowMetaByte, err := json.Marshal(m)
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, flowMetaByte, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (m *Manifest) GetUIBlock(id string) *AppUBLock {
	for i := range m.UIBlocks {
		if m.UIBlocks[i].ID == id {
			return &m.UIBlocks[i]
		}
	}
	return nil
}

func (m *Manifest) GetButton(id string) *UIButton {
	for i := range m.UIButtons {
		if m.UIButtons[i].ID == id {
			return &m.UIButtons[i]
		}
	}
	return nil
}

func (m *Manifest) GetAppConfig(id string) *AppConfig {
	for i := range m.Configs {
		if m.Configs[i].ID == id {
			return &m.Configs[i]
		}
	}
	return nil
}

type AuthResponse struct {
	Status    string `json:"status"`
	ErrorText string `json:"error_text"`
	ErrorCode string `json:"error_code"`
}
