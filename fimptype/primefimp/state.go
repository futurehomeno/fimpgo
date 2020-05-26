package primefimp

import (
	"fmt"
	"reflect"
)

type StateDevice struct {
	Id       int64          `json:"id"`
	Services []StateService `json:"services"`
}

func (sd StateDevice) ContainsService(service string) bool {
	if len(sd.Services) == 0 {
		return false
	}
	for _, ds := range sd.Services {
		if ds.Name == service {
			return true
		}
	}
	return false
}

type StateService struct {
	Name       string           `json:"name"`
	Address    string           `json:"addr"`
	Attributes []StateAttribute `json:"attributes"`
}

func (ss StateService) ContainsAttribute(attribute string) bool {
	if len(ss.Attributes) == 0 {
		return false
	}
	for _, sa := range ss.Attributes {
		if sa.Name == attribute {
			return true
		}
	}
	return false
}

type StateAttribute struct {
	Name   string                `json:"name"`
	Values []StateAttributeValue `json:"values"`
}

type StateAttributeValue struct {
	Timestamp string            `json:"ts"`
	ValType   string            `json:"val_t"`
	Val       interface{}       `json:"val"`
	Props     map[string]string `json:"props"`
}

func (sav StateAttributeValue) GetBoolValue() (bool, error) {
	val, ok := sav.Val.(bool)
	if ok {
		return val, nil
	}
	return false, fmt.Errorf(wrongValueFormat, "bool", reflect.ValueOf(sav.Val))
}

func (sav StateAttributeValue) GetStrArrayValue() ([]string, error) {
	val, ok := sav.Val.([]string)
	if ok {
		return val, nil
	}
	return nil, fmt.Errorf(wrongValueFormat, "[]string", reflect.ValueOf(sav.Val))
}

func (sav StateAttributeValue) GetStrMapValue() (map[string]string, error) {
	val, ok := sav.Val.(map[string]string)
	if ok {
		return val, nil
	}
	return nil, fmt.Errorf(wrongValueFormat, "map[string]string", reflect.ValueOf(sav.Val))
}

func (sav StateAttributeValue) GetIntArrayValue() ([]int64, error) {
	val, ok := sav.Val.([]int64)
	if ok {
		return val, nil
	}
	return nil, fmt.Errorf(wrongValueFormat, "[]int64]", reflect.ValueOf(sav.Val))
}

func (sav StateAttributeValue) GetFloatArrayValue() ([]float64, error) {
	val, ok := sav.Val.([]float64)
	if ok {
		return val, nil
	}
	return nil, fmt.Errorf(wrongValueFormat, "[]float64", reflect.ValueOf(sav.Val))
}

type State struct {
	Devices []StateDevice `json:"devices"`
}

func (s State) FilterDevicesByService(service string) []StateDevice {
	if len(s.Devices) == 0 {
		return nil
	}
	var result []StateDevice

	for _, sd := range s.Devices {
		if sd.ContainsService(service) {
			result = append(result, sd)
			continue
		}

	}
	return result
}

func (s State) FilterDevicesByAttribute(attribute string) []StateDevice {
	if len(s.Devices) == 0 {
		return nil
	}
	var result []StateDevice

	for _, sd := range s.Devices {
		for _, ds := range sd.Services {
			if ds.ContainsAttribute(attribute) {
				result = append(result, sd)
				continue
			}
		}
	}
	return result
}
