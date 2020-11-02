package primefimp

import (
	"encoding/json"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/futurehomeno/fimpgo"
	"github.com/pkg/errors"
	"reflect"
)

/*
	State
*/

type (
	State struct {
		Devices StateDevices `json:"devices"`
	}

	StateDeviceFilter func(*StateDevice) bool

	StateDevice struct {
		Id       int64           `json:"id"`
		Services []*StateService `json:"services"`
	}

	StateService struct {
		Name       string           `json:"name"`
		Address    string           `json:"addr"`
		Attributes []StateAttribute `json:"attributes"`
	}

	StateServiceFilter func(*StateService) bool

	StateDevices []*StateDevice
)

func (s *State) WithFilteredDevices(filter StateDeviceFilter) {
	s.Devices = s.Devices.FilterDevicesByFunc(filter)
}

func (stateDevices StateDevices) FilterDevicesByService(service string) StateDevices {
	if len(stateDevices) == 0 {
		return nil
	}
	var result []*StateDevice

	for _, sd := range stateDevices {
		if sd.ContainsService(service) {
			result = append(result, sd)
			continue
		}

	}
	return result
}

func (stateDevices StateDevices) FilterDevicesByAttribute(attribute string) StateDevices {
	if len(stateDevices) == 0 {
		return nil
	}
	var result []*StateDevice

	for _, sd := range stateDevices {
		for _, ds := range sd.Services {
			if _, ok := ds.FindAttribute(attribute); ok {
				result = append(result, sd)
				continue
			}
		}
	}
	return result
}

func (stateDevices StateDevices) FilterDevicesByFunc(filter StateDeviceFilter) StateDevices {
	var result []*StateDevice
	for _, sd := range stateDevices {
		if filter(sd) {
			result = append(result, sd)
		}
	}
	return result
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
func (sd StateDevice) FilterServices(serviceNames []string) []*StateService {
	var result []*StateService
	for _, sn := range serviceNames {
		for _, stateService := range sd.Services {
			if stateService.Name == sn {
				result = append(result, stateService)
				break
			}
		}
	}
	return result
}
func (sd *StateDevice) WithFilteredServices(filter StateServiceFilter) {
	var temp []*StateService
	for _, ss := range sd.Services {
		if filter(ss) {
			temp = append(temp, ss)
		}
	}
	sd.Services = temp
}

func (ss StateService) FindAttribute(attributeName string) (StateAttribute, bool) {
	for _, serviceAttribute := range ss.Attributes {
		if serviceAttribute.Name == attributeName {
			return serviceAttribute, true
		}
	}
	return StateAttribute{}, false
}

/*
	Attribute
*/

// Note from the documentation: https://github.com/futurehomeno/docs/blob/master/smart-home/core/prime-fimp-states.md
// All attributes will usually have only one value
type StateAttribute struct {
	Name   string                `json:"name"`
	Values []StateAttributeValue `json:"values"`
}

func (sa StateAttribute) validate() error {
	if len(sa.Values) == 0 {
		return fmt.Errorf("empty attribute values for attribute: %+v", sa.Name)
	}
	return nil
}

func (sa StateAttribute) GetFirstValue() (StateAttributeValue, error) {
	if err := sa.validate(); err != nil {
		return StateAttributeValue{}, err
	}
	return sa.Values[0], nil
}

func (sa StateAttribute) GetFirstStringValue() (string, error) {
	if err := sa.validate(); err != nil {
		return "", err
	}
	attrVal := sa.Values[0]
	return attrVal.GetStringValue()
}

func (sa StateAttribute) GetFirstIntValue() (int64, error) {
	if err := sa.validate(); err != nil {
		return -1, err
	}
	attrVal := sa.Values[0]
	return attrVal.GetIntValue()
}

func (sa StateAttribute) GetFirstFloatValue() (float64, error) {
	if err := sa.validate(); err != nil {
		return -1, err
	}
	attrVal := sa.Values[0]
	return attrVal.GetFloatValue()
}

// GetFirstBinaryValue returns the first value of the "binary" attribute
// with the attached timestamp
func (sa StateAttribute) GetFirstBinaryValue() (bool, string, error) {
	if err := sa.validate(); err != nil {
		return false, "", err
	}
	attrVal := sa.Values[0]
	val, err := attrVal.GetBoolValue()
	if err != nil {
		return false, "", err
	}
	return val, attrVal.Timestamp, nil
}

func (sa StateAttribute) GetFirstStrArrayValue() ([]string, error) {
	if err := sa.validate(); err != nil {
		return nil, err
	}
	attrVal := sa.Values[0]
	return attrVal.GetStrArrayValue()
}

func (sa StateAttribute) GetFirstStrMapValue() (map[string]string, error) {
	if err := sa.validate(); err != nil {
		return nil, err
	}
	attrVal := sa.Values[0]
	return attrVal.GetStrMapValue()
}

func (sa StateAttribute) GetFirstIntArrayValue() ([]int64, error) {
	if err := sa.validate(); err != nil {
		return nil, err
	}
	attrVal := sa.Values[0]
	return attrVal.GetIntArrayValue()
}

func (sa StateAttribute) GetFirstIntMapValue() (map[string]int64, error) {
	if err := sa.validate(); err != nil {
		return nil, err
	}
	attrVal := sa.Values[0]
	return attrVal.GetIntMapValue()
}

func (sa StateAttribute) GetFirstFloatArrayValue() ([]float64, error) {
	if err := sa.validate(); err != nil {
		return nil, err
	}
	attrVal := sa.Values[0]
	return attrVal.GetFloatArrayValue()
}

func (sa StateAttribute) GetFirstFloatMapValue() (map[string]float64, error) {
	if err := sa.validate(); err != nil {
		return nil, err
	}
	attrVal := sa.Values[0]
	return attrVal.GetFloatMapValue()
}

func (sa StateAttribute) GetFirstBoolMapValue() (map[string]bool, error) {
	if err := sa.validate(); err != nil {
		return nil, err
	}
	attrVal := sa.Values[0]
	return attrVal.GetBoolMapValue()
}

func (sa StateAttribute) GetFirstPropAsString(propName string) (string, error) {
	if err := sa.validate(); err != nil {
		return "", err
	}
	props := sa.Values[0].Props
	propVal, ok := props[propName]
	if !ok {
		return "", fmt.Errorf("cannot find attribute property: %+v", propName)
	}
	return propVal, nil
}

/*
	Attribute Value
*/
type StateAttributeValue struct {
	Timestamp string            `json:"ts"`
	ValType   string            `json:"val_t"`
	Val       interface{}       `json:"val"`
	Props     map[string]string `json:"props"`
}

func (sav *StateAttributeValue) parse() error {
	b, err := json.Marshal(sav)
	if err != nil {
		return errors.Wrap(err, "marshalling")
	}
	switch sav.ValType {
	case fimpgo.VTypeString:
		if sav.Val, err = jsonparser.GetString(b, fimpgo.Val); err != nil {
			return err
		}
	case fimpgo.VTypeBool:
		if sav.Val, err = jsonparser.GetBoolean(b, fimpgo.Val); err != nil {
			return err
		}
	case fimpgo.VTypeInt:
		if sav.Val, err = jsonparser.GetInt(b, fimpgo.Val); err != nil {
			return err
		}
	case fimpgo.VTypeFloat:
		if sav.Val, err = jsonparser.GetFloat(b, fimpgo.Val); err != nil {
			return err
		}
	case fimpgo.VTypeBoolArray:
		var val []bool
		_, err = jsonparser.ArrayEach(b, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			item, _ := jsonparser.ParseBoolean(value)
			val = append(val, item)
		}, fimpgo.Val)
		if err != nil {
			return err
		}
		sav.Val = val
	case fimpgo.VTypeStrArray:
		var val []string
		if _, err = jsonparser.ArrayEach(b, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			item, _ := jsonparser.ParseString(value)
			val = append(val, item)
		}, fimpgo.Val); err != nil {
			return err
		}
		sav.Val = val
	case fimpgo.VTypeIntArray:
		var val []int64
		if _, err = jsonparser.ArrayEach(b, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			item, _ := jsonparser.ParseInt(value)
			val = append(val, item)
		}, fimpgo.Val); err != nil {
			return err
		}
		sav.Val = val
	case fimpgo.VTypeFloatArray:
		var val []float64
		if _, err = jsonparser.ArrayEach(b, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			item, _ := jsonparser.ParseFloat(value)
			val = append(val, item)
		}, fimpgo.Val); err != nil {
			return nil
		}
		sav.Val = val
	case fimpgo.VTypeStrMap:
		val := make(map[string]string)
		if err = jsonparser.ObjectEach(b, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			val[string(key)], err = jsonparser.ParseString(value)
			return nil
		}, fimpgo.Val); err != nil {
			return err
		}
		sav.Val = val
	case fimpgo.VTypeIntMap:
		val := make(map[string]int64)
		if err = jsonparser.ObjectEach(b, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			val[string(key)], err = jsonparser.ParseInt(value)
			return nil
		}, fimpgo.Val); err != nil {
			return err
		}
		sav.Val = val
	case fimpgo.VTypeFloatMap:
		val := make(map[string]bool)
		if err = jsonparser.ObjectEach(b, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			val[string(key)], err = jsonparser.ParseBoolean(value)
			return nil
		}, fimpgo.Val); err != nil {
			return err
		}
		sav.Val = val
	case fimpgo.VTypeBoolMap:
		val := make(map[string]bool)
		if err = jsonparser.ObjectEach(b, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			val[string(key)], err = jsonparser.ParseBoolean(value)
			return nil
		}, fimpgo.Val); err != nil {
			return err
		}
		sav.Val = val
	}
	return nil
}

func (sav StateAttributeValue) GetStringValue() (string, error) {
	if err := (&sav).parse(); err != nil {
		return "", errors.Wrap(err, "parsing")
	}
	val, ok := sav.Val.(string)
	if ok {
		return val, nil
	}
	return "", fmt.Errorf(wrongValueFormat, "string", reflect.ValueOf(sav.Val))
}

func (sav StateAttributeValue) GetIntValue() (int64, error) {
	if err := (&sav).parse(); err != nil {
		return -1, errors.Wrap(err, "parsing")
	}
	val, ok := sav.Val.(int64)
	if ok {
		return val, nil
	}
	return -1, fmt.Errorf(wrongValueFormat, "int64", reflect.ValueOf(sav.Val))
}

func (sav StateAttributeValue) GetFloatValue() (float64, error) {
	if err := (&sav).parse(); err != nil {
		return -1, errors.Wrap(err, "parsing")
	}
	val, ok := sav.Val.(float64)
	if ok {
		return val, nil
	}
	return -1, fmt.Errorf(wrongValueFormat, "float64", reflect.ValueOf(sav.Val))
}

func (sav StateAttributeValue) GetBoolValue() (bool, error) {
	if err := (&sav).parse(); err != nil {
		return false, errors.Wrap(err, "parsing")
	}
	val, ok := sav.Val.(bool)
	if ok {
		return val, nil
	}
	return false, fmt.Errorf(wrongValueFormat, "bool", reflect.ValueOf(sav.Val))
}

func (sav StateAttributeValue) GetStrArrayValue() ([]string, error) {
	if err := (&sav).parse(); err != nil {
		return nil, errors.Wrap(err, "parsing")
	}
	val, ok := sav.Val.([]string)
	if ok {
		return val, nil
	}
	return nil, fmt.Errorf(wrongValueFormat, "[]string", reflect.ValueOf(sav.Val))
}

func (sav StateAttributeValue) GetStrMapValue() (map[string]string, error) {
	if err := (&sav).parse(); err != nil {
		return nil, errors.Wrap(err, "parsing")
	}
	strMapVal, ok := sav.Val.(map[string]string)
	if ok {
		return strMapVal, nil
	}

	iMapVal, ok := sav.Val.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf(wrongValueFormat, "map[string]string", reflect.ValueOf(sav.Val))
	}

	result := map[string]string{}
	for k, v := range iMapVal {
		result[k] = fmt.Sprint(v)
	}
	return result, nil

}

func (sav StateAttributeValue) GetIntArrayValue() ([]int64, error) {
	if err := (&sav).parse(); err != nil {
		return nil, errors.Wrap(err, "parsing")
	}
	val, ok := sav.Val.([]int64)
	if ok {
		return val, nil
	}
	return nil, fmt.Errorf(wrongValueFormat, "[]int64]", reflect.ValueOf(sav.Val))
}

func (sav StateAttributeValue) GetFloatArrayValue() ([]float64, error) {
	if err := (&sav).parse(); err != nil {
		return nil, errors.Wrap(err, "parsing")
	}
	val, ok := sav.Val.([]float64)
	if ok {
		return val, nil
	}
	return nil, fmt.Errorf(wrongValueFormat, "[]float64", reflect.ValueOf(sav.Val))
}

func (sav StateAttributeValue) GetFloatMapValue() (map[string]float64, error) {
	if err := (&sav).parse(); err != nil {
		return nil, errors.Wrap(err, "parsing")
	}
	val, ok := sav.Val.(map[string]float64)
	if ok {
		return val, nil
	}
	return nil, fmt.Errorf(wrongValueFormat, "map[string]float64", reflect.ValueOf(sav.Val))
}

func (sav StateAttributeValue) GetIntMapValue() (map[string]int64, error) {
	if err := (&sav).parse(); err != nil {
		return nil, errors.Wrap(err, "parsing")
	}
	val, ok := sav.Val.(map[string]int64)
	if ok {
		return val, nil
	}
	return nil, fmt.Errorf(wrongValueFormat, "map[string]int64", reflect.ValueOf(sav.Val))
}

func (sav StateAttributeValue) GetBoolMapValue() (map[string]bool, error) {
	if err := (&sav).parse(); err != nil {
		return nil, errors.Wrap(err, "parsing")
	}
	val, ok := sav.Val.(map[string]bool)
	if ok {
		return val, nil
	}
	return nil, fmt.Errorf(wrongValueFormat, "map[string]bool", reflect.ValueOf(sav.Val))
}
