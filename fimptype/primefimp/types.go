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

	CmdGet    = "get"
	CmdSet    = "set"
	CmdEdit   = "edit"
	CmdDelete = "delete"
)

// Top level element for commands
type Request struct {
	Cmd       string       `json:"cmd"`
	Component interface{}  `json:"component"`
	Param     RequestParam `json:"param"`
	RequestID interface{}  `json:"requestId,omitempty"`
	Id        interface{}  `json:"id,omitempty"`
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
	Name string `json:"name"`
}

type Device struct {
	Fimp          Fimp                   `json:"fimp"`
	Client        Client                 `json:"client"`
	Functionality string                 `json:"functionality"`
	Service       map[string]Service     `json:"services"`
	ID            int                    `json:"id"`
	Lrn           bool                   `json:"lrn"`
	Model         string                 `json:"model"`
	Param         map[string]interface{} `json:"param"`
	Problem       bool                   `json:"problem"`
	Room          int                    `json:"room"`
	Changes       map[string]interface{} `json:"changes"`
	ThingID       int                    `json:"thing"`
}

type Thing struct {
	ID      int               `json:"id"`
	Address string            `json:"addr"`
	Name    string            `json:"name"`
	Devices []int             `json:"devices,omitempty"`
	Props   map[string]string `json:"props,omitempty"`
	RoomID  int               `json:"room"`
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
	Type    string     `json:"type"`
	Area    int        `json:"area,omitempty"`
	Outside bool       `json:"outside"`
}

type RoomParams struct {
	Heating RoomHeating `json:"heating"`
}

type RoomHeating struct {
	Desired float64 `json:"desired"`
	Target  float64 `json:"target"`
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
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
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
	Device map[int]ActionDevice `json:"device"`
	Room   map[int]ActionRoom   `json:"room"`
}

type Mode struct {
	Id     string     `json:"id"`
	Action ModeAction `json:"action"`
}
