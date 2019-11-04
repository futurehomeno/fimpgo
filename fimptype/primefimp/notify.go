package primefimp

import (
	"encoding/json"
	"errors"

	"github.com/futurehomeno/fimpgo"
)

type Notify struct {
	Errors     interface{}     `json:"errors"`
	Cmd        string          `json:"cmd"`
	Component  string          `json:"component"`
	ParamRaw   json.RawMessage `json:"param"`
	ChangesRaw json.RawMessage `json:"changes"`
	Success    bool            `json:"success"`
	Id         interface{}     `json:"id,omitempty"`
}

type DeleteChange struct {
	ID int `json:"id"`
}

func FimpToNotify(msg *fimpgo.Message) (*Notify, error) {
	if msg.Payload.Type != "evt.pd7.notify" {
		return nil, errors.New("wrong fimp msg type")
	}
	notify := Notify{}
	err := msg.Payload.GetObjectValue(&notify)
	if err != nil {
		return nil, err
	}

	return &notify, err
}

func (ntf *Notify) GetDevice() *Device {
	if ntf.Component == ComponentDevice {
		var result Device
		err := json.Unmarshal(ntf.ParamRaw, &result)
		if err != nil {
			return nil
		}
		return &result
	}
	return nil
}

func (ntf *Notify) GetThing() *Thing {
	if ntf.Component == ComponentThing {
		var result Thing
		err := json.Unmarshal(ntf.ParamRaw, &result)
		if err != nil {
			return nil
		}
		return &result
	}
	return nil
}

func (ntf *Notify) GetRoom() *Room {
	if ntf.Component == ComponentRoom {
		var result Room
		err := json.Unmarshal(ntf.ParamRaw, &result)
		if err != nil {
			return nil
		}
		return &result
	}
	return nil
}

func (ntf *Notify) GetArea() *Area {
	if ntf.Component == ComponentArea {
		var result Area
		err := json.Unmarshal(ntf.ParamRaw, &result)
		if err != nil {
			return nil
		}
		return &result
	}
	return nil
}

func (ntf *Notify) GetHouse() *House {
	if ntf.Component == ComponentHouse {
		var result House
		err := json.Unmarshal(ntf.ParamRaw, &result)
		if err != nil {
			return nil
		}
		return &result
	}
	return nil
}

func (ntf *Notify) GetShortcut() *Shortcut {
	if ntf.Component == ComponentShortcut {
		var result Shortcut
		err := json.Unmarshal(ntf.ParamRaw, &result)
		if err != nil {
			return nil
		}
		return &result
	}
	return nil
}

func (ntf *Notify) GetTimer() *Timer {
	if ntf.Component == ComponentTimer {
		var result Timer
		err := json.Unmarshal(ntf.ParamRaw, &result)
		if err != nil {
			return nil
		}
		return &result
	}
	return nil
}

func (ntf *Notify) GetHub() *Hub {
	if ntf.Component == ComponentHub {
		var result Hub
		err := json.Unmarshal(ntf.ParamRaw, &result)
		if err != nil {
			return nil
		}
		return &result
	}
	return nil
}

func (ntf *Notify) GetDeleteChange() *DeleteChange {
	var result DeleteChange
	err := json.Unmarshal(ntf.ChangesRaw, &result)
	if err != nil {
		return nil
	}
	return &result
}
