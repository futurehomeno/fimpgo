package primefimp

import (
	"time"
)

const (
	ComponentDevice   = "device"
	ComponentThing    = "thing"
	ComponentRoom     = "room"
	ComponentArea     = "area"
	ComponentHouse    = "house"
	ComponentHub      = "hub"
	ComponentShortcut = "shortcut"
	ComponentMode     = "mode"
	ComponentTimer    = "timer"
	ComponentService  = "service"
	ComponentState    = "state"

	CmdGet    = "get"
	CmdSet    = "set"
	CmdEdit   = "edit"
	CmdDelete = "delete"
	CmdAdd    = "add"

	wrongValueFormat = "wrong value type. expected %+v, got %+v"

	ServiceOutBinSwitch = "out_bin_switch"
	ServiceOutLvlSwitch = "out_lvl_switch"
	ServiceThermostat   = "thermostat"
	ServiceColorControl = "color_ctrl"
	ServiceBattery      = "battery"

	// sensors
	ServiceSensorTemp    = "sensor_temp"
	ServiceSensorContact = "sensor_contact"
	ServiceSensorLumin   = "sensor_lumin"
	ServiceSensorHumid   = "sensor_humid"
)

// Top level element for commands
type Request struct {
	Cmd       string        `json:"cmd"`
	Component interface{}   `json:"component"`
	Param     *RequestParam `json:"param,omitempty"`
	RequestID interface{}   `json:"requestId,omitempty"`
	Id        interface{}   `json:"id,omitempty"`
}

type RequestParam struct {
	Id         int      `json:"id,omitempty"`
	Components []string `json:"components,omitempty"`
}

type Fimp struct {
	Adapter string `json:"adapter"`
	Address string `json:"address"`
	Group   string `json:"group"`
}

type Client struct {
	Name          *string `json:"name,omitempty"`
	OpenStateType *string `json:"openStateType,omitempty"`
}

type Device struct {
	Fimp          Fimp                   `json:"fimp"`
	Client        Client                 `json:"client"`
	Functionality *string                `json:"functionality"`
	Service       map[string]Service     `json:"services"`
	ID            int                    `json:"id"`
	Lrn           bool                   `json:"lrn"`
	Model         string                 `json:"model"`
	ModelAlias    string                 `json:"modelAlias"`
	Param         map[string]interface{} `json:"param"`
	Problem       bool                   `json:"problem"`
	Room          *int                   `json:"room"`
	Changes       map[string]interface{} `json:"changes"`
	ThingID       *int                   `json:"thing"`
	Type          map[string]interface{} `json:"type"`
}

type Thing struct {
	ID      int                    `json:"id"`
	Address string                 `json:"addr"`
	Name    string                 `json:"name"`
	Devices []int                  `json:"devices,omitempty"`
	Props   map[string]interface{} `json:"props,omitempty"`
	RoomID  int                    `json:"room"`
}

type House struct {
	Learning interface{} `json:"learning"`
	Mode     string      `json:"mode"`
	Time     time.Time   `json:"time"`
}

type Room struct {
	Alias   string     `json:"alias"`
	ID      int        `json:"id"`
	Param   RoomParams `json:"param"`
	Client  Client     `json:"client"`
	Type    *string    `json:"type"`
	Area    *int       `json:"area"`
	Outside bool       `json:"outside"`
}

type RoomParams struct {
	Heating  RoomHeating `json:"heating"`
	Lighting interface{} `json:"lighting"`
	Security interface{} `json:"security"`
	Sensors  []string    `json:"sensors"`
	Shading  interface{} `json:"shading"`
	Triggers interface{} `json:"triggers"`
}

type RoomHeating struct {
	Desired    float64 `json:"desired"`
	Target     float64 `json:"target"`
	Thermostat bool    `json:"thermostat"`
	Actuator   bool    `json:"actuator"`
	Power      string  `json:"power"`
}

type Service struct {
	Addr       string                 `json:"addr,omitempty"`
	Enabled    bool                   `json:"enabled,omitempty"`
	Interfaces []string               `json:"intf"`
	Props      map[string]interface{} `json:"props"`
}

type UserInfo struct {
	UID  string   `json:"uuid,omitempty"`
	Name UserName `json:"name,omitempty"`
}

type UserName struct {
	Fullname string `json:"fullname,omitempty"`
}

type Area struct {
	ID    int       `json:"id"`
	Name  string    `json:"name"`
	Type  string    `json:"type"`
	Props AreaProps `json:"props"`
}

type AreaProps struct {
	HNumber string `json:"hNumber"`
	TransNr string `json:"transNr"`
}

type ActionDevice map[string]interface{}

type ActionRoom map[string]interface{}

type ShortcutAction struct {
	Device map[int]ActionDevice `json:"device"`
	Room   map[int]ActionRoom   `json:"room"`
}

type Shortcut struct {
	ID     int            `json:"id"`
	Client Client         `json:"client"`
	Action ShortcutAction `json:"action"`
}

type HubMode struct {
	Current  string `json:"current"`
	Previous string `json:"prev"`
}

type Hub struct {
	Mode HubMode `json:"mode"`
}

type ModeAction struct {
	Device ActionDevice `json:"device"`
	Room   ActionRoom   `json:"room"`
}

type Mode struct {
	Id     string     `json:"id"`
	Action ModeAction `json:"action"`
}

type TimerAction struct {
	Type     string
	Shortcut int
	Mode     string
	Action   ShortcutAction
}

type Timer struct {
	Action  TimerAction
	Client  Client                 `json:"client"`
	Enabled bool                   `json:"enabled"`
	Time    map[string]interface{} `json:"time"`
	ID      int                    `json:"id"`
}

type VincServices struct {
	FireAlarm map[string]interface{} `json:"fireAlarm"`
}
