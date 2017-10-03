package fimpgo

import "time"
import (
	"encoding/json"
	"errors"
	"github.com/buger/jsonparser"
	"github.com/satori/go.uuid"
)

const (
	VTypeString     = "string"
	VTypeInt        = "int"
	VTypeFloat      = "float"
	VTypeBool       = "bool"
	VTypeStrMap     = "str_map"
	VTypeIntMap     = "int_map"
	VTypeFloatMap   = "float_map"
	VTypeBoolMap    = "bool_map"
	VTypeStrArray   = "str_array"
	VTypeIntArray   = "int_array"
	VTypeFloatArray = "float_array"
	VTypeBoolArray  = "bool_array"
	VTypeObject     = "object"
	VTypeNull       = "null"
)

type Props map[string]string
type Tags []string

type FimpMessage struct {
	Type          string      `json:"type"`
	Service       string      `json:"serv"`
	ValueType     string      `json:"val_t"`
	Value         interface{} `json:"val"`
	ValueObj      []byte      `json:"-"`
	Tags          Tags        `json:"tags"`
	Properties    Props       `json:"props"`
	Version       string      `json:"ver"`
	CorrelationID string      `json:"corid"`
	CreationTime  time.Time   `json:"ctime"`
	UID           string      `json:"uid"`
}

func (msg *FimpMessage) SetValue(value interface{}, valType string) {
	msg.Value = value
	msg.ValueType = valType
}

func (msg *FimpMessage) GetIntValue() (int64, error) {
	val, ok := msg.Value.(int64)
	if ok {
		return val, nil
	}
	return 0, errors.New("Wrong value type")
}

func (msg *FimpMessage) GetStringValue() (string, error) {
	val, ok := msg.Value.(string)
	if ok {
		return val, nil
	}
	return "", errors.New("Wrong value type")
}

func (msg *FimpMessage) GetBoolValue() (bool, error) {
	val, ok := msg.Value.(bool)
	if ok {
		return val, nil
	}
	return false, errors.New("Wrong value type")
}

func (msg *FimpMessage) GetFloatValue() (float64, error) {
	val, ok := msg.Value.(float64)
	if ok {
		return val, nil
	}
	return 0, errors.New("Wrong value type")
}

func (msg *FimpMessage) GetStrArrayValue() ([]string, error) {
	val, ok := msg.Value.([]string)
	if ok {
		return val, nil
	}
	return nil, errors.New("Wrong value type")
}

func (msg *FimpMessage) GetIntArrayValue() ([]int64, error) {
	val, ok := msg.Value.([]int64)
	if ok {
		return val, nil
	}
	return nil, errors.New("Wrong value type")
}

func (msg *FimpMessage) GetFloatArrayValue() ([]float64, error) {
	val, ok := msg.Value.([]float64)
	if ok {
		return val, nil
	}
	return nil, errors.New("Wrong value type")
}

func (msg *FimpMessage) GetBoolArrayValue() ([]bool, error) {
	val, ok := msg.Value.([]bool)
	if ok {
		return val, nil
	}
	return nil, errors.New("Wrong value type")
}

func (msg *FimpMessage) GetStrMapValue() (map[string]string, error) {
	val, ok := msg.Value.(map[string]string)
	if ok {
		return val, nil
	}
	return nil, errors.New("Wrong value type")
}

func (msg *FimpMessage) GetIntMapValue() (map[string]int64, error) {
	val, ok := msg.Value.(map[string]int64)
	if ok {
		return val, nil
	}
	return nil, errors.New("Wrong value type")
}

func (msg *FimpMessage) GetFloatMapValue() (map[string]float64, error) {
	val, ok := msg.Value.(map[string]float64)
	if ok {
		return val, nil
	}
	return nil, errors.New("Wrong value type")
}

func (msg *FimpMessage) GetBoolMapValue() (map[string]bool, error) {
	val, ok := msg.Value.(map[string]bool)
	if ok {
		return val, nil
	}
	return nil, errors.New("Wrong value type")
}

func (msg *FimpMessage) GetRawObjectValue() []byte {
	return msg.ValueObj
}

func (msg *FimpMessage) GetObjectValue(objectBindVar interface{}) error {
	return json.Unmarshal(msg.ValueObj, objectBindVar)
}

func (msg *FimpMessage) SerializeToJson() ([]byte, error) {
	jsonBA, err := json.Marshal(msg)
	return jsonBA, err

}



func NewMessage(type_ string, service string, valueType string, value interface{}, props Props, tags Tags, requestMessage *FimpMessage) *FimpMessage {
	msg := FimpMessage{Type: type_,
		Service:    service,
		ValueType:  valueType,
		Value:      value,
		Tags:       tags,
		Properties: props,
		Version:    "1",
		CreationTime:time.Now(),
		UID:uuid.NewV4().String(),
	}

	if requestMessage != nil {
		msg.CorrelationID = requestMessage.UID
	}

	return &msg
}

func NewNullMessage(type_ string, service string, props Props, tags Tags, requestMessage *FimpMessage) *FimpMessage {
	return NewMessage(type_, service, VTypeNull, nil, props, tags, requestMessage)
}

func NewStringMessage(type_ string, service string, value string, props Props, tags Tags, requestMessage *FimpMessage) *FimpMessage {
	return NewMessage(type_, service, VTypeString, value, props, tags, requestMessage)
}

func NewIntMessage(type_ string, service string, value int64, props Props, tags Tags, requestMessage *FimpMessage) *FimpMessage {
	return NewMessage(type_, service, VTypeInt, value, props, tags, requestMessage)
}

func NewFloatMessage(type_ string, service string, value float64, props Props, tags Tags, requestMessage *FimpMessage) *FimpMessage {
	return NewMessage(type_, service, VTypeFloat, value, props, tags, requestMessage)
}

func NewBoolMessage(type_ string, service string, value bool, props Props, tags Tags, requestMessage *FimpMessage) *FimpMessage {
	return NewMessage(type_, service, VTypeBool, value, props, tags, requestMessage)
}

func NewStrArrayMessage(type_ string, service string, value []string, props Props, tags Tags, requestMessage *FimpMessage) *FimpMessage {
	return NewMessage(type_, service, VTypeStrArray, value, props, tags, requestMessage)
}

func NewIntArrayMessage(type_ string, service string, value []int64, props Props, tags Tags, requestMessage *FimpMessage) *FimpMessage {
	return NewMessage(type_, service, VTypeIntArray, value, props, tags, requestMessage)
}

func NewFloatArrayMessage(type_ string, service string, value []float64, props Props, tags Tags, requestMessage *FimpMessage) *FimpMessage {
	return NewMessage(type_, service, VTypeFloatArray, value, props, tags, requestMessage)
}

func NewBoolArrayMessage(type_ string, service string, value []bool, props Props, tags Tags, requestMessage *FimpMessage) *FimpMessage {
	return NewMessage(type_, service, VTypeBoolArray, value, props, tags, requestMessage)
}

func NewStrMapMessage(type_ string, service string, value map[string]string, props Props, tags Tags, requestMessage *FimpMessage) *FimpMessage {
	return NewMessage(type_, service, VTypeStrMap, value, props, tags, requestMessage)
}

func NewIntMapMessage(type_ string, service string, value map[string]int64, props Props, tags Tags, requestMessage *FimpMessage) *FimpMessage {
	return NewMessage(type_, service, VTypeIntMap, value, props, tags, requestMessage)
}

func NewFloatMapMessage(type_ string, service string, value map[string]float64, props Props, tags Tags, requestMessage *FimpMessage) *FimpMessage {
	return NewMessage(type_, service, VTypeFloatMap, value, props, tags, requestMessage)
}

func NewBoolMapMessage(type_ string, service string, value map[string]bool, props Props, tags Tags, requestMessage *FimpMessage) *FimpMessage {
	return NewMessage(type_, service, VTypeFloatMap, value, props, tags, requestMessage)
}

func NewMessageFromBytes(msg []byte) (*FimpMessage, error) {
	fimpmsg := FimpMessage{}
	var err error
	fimpmsg.Type, err = jsonparser.GetString(msg, "type")
	fimpmsg.Service, err = jsonparser.GetString(msg, "serv")
	fimpmsg.ValueType, err = jsonparser.GetString(msg, "val_t")
	fimpmsg.UID, err = jsonparser.GetString(msg, "uid")
	fimpmsg.CorrelationID, err = jsonparser.GetString(msg, "corid")
	switch fimpmsg.ValueType {
	case VTypeString:
		fimpmsg.Value, err = jsonparser.GetString(msg, "val")
	case VTypeBool:
		fimpmsg.Value, err = jsonparser.GetBoolean(msg, "val")
	case VTypeInt:
		fimpmsg.Value, err = jsonparser.GetInt(msg, "val")
	case VTypeFloat:
		fimpmsg.Value, err = jsonparser.GetFloat(msg, "val")
	case VTypeBoolArray:
		val := []bool{}
		jsonparser.ArrayEach(msg, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			item, _ := jsonparser.ParseBoolean(value)
			val = append(val, item)
		}, "val")

		fimpmsg.Value = val
	case VTypeStrArray:
		val := []string{}
		jsonparser.ArrayEach(msg, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			item, _ := jsonparser.ParseString(value)
			val = append(val, item)
		}, "val")

		fimpmsg.Value = val
	case VTypeIntArray:
		val := []int64{}
		jsonparser.ArrayEach(msg, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			item, _ := jsonparser.ParseInt(value)
			val = append(val, item)
		}, "val")
		fimpmsg.Value = val
	case VTypeFloatArray:
		val := []float64{}
		jsonparser.ArrayEach(msg, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			item, _ := jsonparser.ParseFloat(value)
			val = append(val, item)
		}, "val")
		fimpmsg.Value = val

	case VTypeStrMap:
		val := make(map[string]string)
		jsonparser.ObjectEach(msg, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			val[string(key)],err = jsonparser.ParseString(value)
			return nil
		}, "val")
		fimpmsg.Value = val

	case VTypeIntMap:
		val := make(map[string]int64)
		jsonparser.ObjectEach(msg, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			val[string(key)], err = jsonparser.ParseInt(value)
			return nil
		}, "val")
		fimpmsg.Value = val

	case VTypeFloatMap:
		val := make(map[string]float64)
		jsonparser.ObjectEach(msg, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			val[string(key)], err = jsonparser.ParseFloat(value)
			return nil
		}, "val")
		fimpmsg.Value = val

	case VTypeBoolMap:
		val := make(map[string]bool)
		jsonparser.ObjectEach(msg, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			val[string(key)], err = jsonparser.ParseBoolean(value)
			return nil
		}, "val")
		fimpmsg.Value = val

	case VTypeObject:
		fimpmsg.ValueObj, _, _, err = jsonparser.Get(msg, "val")

	}

	return &fimpmsg, err

}
