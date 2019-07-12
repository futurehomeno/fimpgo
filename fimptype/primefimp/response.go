package primefimp

import (
	"encoding/json"
	"errors"
	"github.com/futurehomeno/fimpgo"
)

type Response struct {
	Errors    interface{}                `json:"errors"`
	Cmd       string                     `json:"cmd"`
	ParamRaw  map[string]json.RawMessage `json:"param"`
	RequestID interface{}                `json:",requestId"`
	Success   bool                       `json:"success"`
	Id        interface{}                `json:"id,omitempty"`
}

func FimpToResponse(msg *fimpgo.FimpMessage) (*Response, error) {
	if msg.Type != "evt.pd7.response" {
		return nil, errors.New("wrong fimp msg type")
	}
	response := Response{}
	err := msg.GetObjectValue(&response)
	if err != nil {
		return nil, err
	}
	return &response, err
}

func (resp *Response) GetDevices() []Device {
	param, ok := resp.ParamRaw[ComponentDevice]
	if !ok {
		return nil
	}
	var result []Device
	err := json.Unmarshal(param, &result)
	if err != nil {
		return nil
	}
	return result
}

func (resp *Response) GetRooms() []Room {
	param, ok := resp.ParamRaw[ComponentRoom]
	if !ok {
		return nil
	}
	var result []Room
	err := json.Unmarshal(param, &result)
	if err != nil {
		return nil
	}
	return result
}

func (resp *Response) GetThings() []Thing {
	param, ok := resp.ParamRaw[ComponentThing]
	if !ok {
		return nil
	}
	var result []Thing
	err := json.Unmarshal(param, &result)
	if err != nil {
		return nil
	}
	return result
}

func (resp *Response) GetAreas() []Area {
	param, ok := resp.ParamRaw[ComponentArea]
	if !ok {
		return nil
	}
	var result []Area
	err := json.Unmarshal(param, &result)
	if err != nil {
		return nil
	}
	return result
}

func (resp *Response) GetHouse() *House {
	param, ok := resp.ParamRaw[ComponentHouse]
	if !ok {
		return nil
	}
	var result House
	err := json.Unmarshal(param, &result)
	if err != nil {
		return nil
	}
	return &result
}

func (resp *Response) GetShortcuts() []Shortcut {
	param, ok := resp.ParamRaw[ComponentShortcut]
	if !ok {
		return nil
	}
	var result []Shortcut
	err := json.Unmarshal(param, &result)
	if err != nil {
		return nil
	}
	return result
}

func (resp *Response) GetModes() []Mode {
	param, ok := resp.ParamRaw[ComponentMode]
	if !ok {
		return nil
	}
	var result []Mode
	err := json.Unmarshal(param, &result)
	if err != nil {
		return nil
	}
	return result
}

func (resp *Response) GetTimers() []Mode {
	param, ok := resp.ParamRaw[ComponentTimer]
	if !ok {
		return nil
	}
	var result []Mode
	err := json.Unmarshal(param, &result)
	if err != nil {
		return nil
	}
	return result
}
