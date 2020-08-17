package fimpgo

import (
	"encoding/base64"
	"fmt"
	"reflect"
	"time"
)
import (
	"encoding/json"
	"github.com/buger/jsonparser"
	"github.com/satori/go.uuid"
)

const (
	TimeFormat       = "2006-01-02T15:04:05.999Z07:00"
	VTypeString      = "string"
	VTypeInt         = "int"
	VTypeFloat       = "float"
	VTypeBool        = "bool"
	VTypeStrMap      = "str_map"
	VTypeIntMap      = "int_map"
	VTypeFloatMap    = "float_map"
	VTypeBoolMap     = "bool_map"
	VTypeStrArray    = "str_array"
	VTypeIntArray    = "int_array"
	VTypeFloatArray  = "float_array"
	VTypeBoolArray   = "bool_array"
	VTypeObject      = "object"
	VTypeBase64      = "base64"
	VTypeBinary      = "bin"
	VTypeNull        = "null"
	wrongValueFormat = "wrong value type. expected %+v, got %+v"

	Val = "val"
)

type Props map[string]string
type Tags []string

type FimpMessage struct {
	Type            string      `json:"type"`
	Service         string      `json:"serv"`
	ValueType       string      `json:"val_t"`
	Value           interface{} `json:"val"`
	ValueObj        []byte      `json:"-"`
	Tags            Tags        `json:"tags"`
	Properties      Props       `json:"props"`
	Version         string      `json:"ver"`
	CorrelationID   string      `json:"corid"`
	ResponseToTopic string      `json:"resp_to,omitempty"`
	Source          string      `json:"src,omitempty"`
	CreationTime    string      `json:"ctime"`
	UID             string      `json:"uid"`
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
	return 0, fmt.Errorf(wrongValueFormat, "int64", reflect.ValueOf(msg.Value))
}

func (msg *FimpMessage) GetStringValue() (string, error) {
	val, ok := msg.Value.(string)
	if ok {
		return val, nil
	}
	return "", fmt.Errorf(wrongValueFormat, "string", reflect.ValueOf(msg.Value))
}

func (msg *FimpMessage) GetBoolValue() (bool, error) {
	val, ok := msg.Value.(bool)
	if ok {
		return val, nil
	}
	return false, fmt.Errorf(wrongValueFormat, "bool", reflect.ValueOf(msg.Value))
}

func (msg *FimpMessage) GetFloatValue() (float64, error) {
	val, ok := msg.Value.(float64)
	if ok {
		return val, nil
	}
	return 0, fmt.Errorf(wrongValueFormat, "float64", reflect.ValueOf(msg.Value))
}

func (msg *FimpMessage) GetStrArrayValue() ([]string, error) {
	val, ok := msg.Value.([]string)
	if ok {
		return val, nil
	}
	return nil, fmt.Errorf(wrongValueFormat, "[]string", reflect.ValueOf(msg.Value))
}

func (msg *FimpMessage) GetIntArrayValue() ([]int64, error) {
	val, ok := msg.Value.([]int64)
	if ok {
		return val, nil
	}
	return nil, fmt.Errorf(wrongValueFormat, "[]int64]", reflect.ValueOf(msg.Value))
}

func (msg *FimpMessage) GetFloatArrayValue() ([]float64, error) {
	val, ok := msg.Value.([]float64)
	if ok {
		return val, nil
	}
	return nil, fmt.Errorf(wrongValueFormat, "[]float64", reflect.ValueOf(msg.Value))
}

func (msg *FimpMessage) GetBoolArrayValue() ([]bool, error) {
	val, ok := msg.Value.([]bool)
	if ok {
		return val, nil
	}
	return nil, fmt.Errorf(wrongValueFormat, "[]bool", reflect.ValueOf(msg.Value))
}

func (msg *FimpMessage) GetStrMapValue() (map[string]string, error) {
	val, ok := msg.Value.(map[string]string)
	if ok {
		return val, nil
	}
	return nil, fmt.Errorf(wrongValueFormat, "map[string]string", reflect.ValueOf(msg.Value))
}

func (msg *FimpMessage) GetIntMapValue() (map[string]int64, error) {
	val, ok := msg.Value.(map[string]int64)
	if ok {
		return val, nil
	}
	return nil, fmt.Errorf(wrongValueFormat, "map[string]int64", reflect.ValueOf(msg.Value))
}

func (msg *FimpMessage) GetFloatMapValue() (map[string]float64, error) {
	val, ok := msg.Value.(map[string]float64)
	if ok {
		return val, nil
	}
	return nil, fmt.Errorf(wrongValueFormat, "map[string]float64", reflect.ValueOf(msg.Value))
}

func (msg *FimpMessage) GetBoolMapValue() (map[string]bool, error) {
	val, ok := msg.Value.(map[string]bool)
	if ok {
		return val, nil
	}
	return nil, fmt.Errorf(wrongValueFormat, "map[string]bool", reflect.ValueOf(msg.Value))
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
		Service:      service,
		ValueType:    valueType,
		Value:        value,
		Tags:         tags,
		Properties:   props,
		Version:      "1",
		CreationTime: time.Now().Format(TimeFormat),
		UID:          uuid.NewV4().String(),
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
	return NewMessage(type_, service, VTypeBoolMap, value, props, tags, requestMessage)
}

func NewObjectMessage(type_ string, service string, value interface{}, props Props, tags Tags, requestMessage *FimpMessage) *FimpMessage {
	return NewMessage(type_, service, VTypeObject, value, props, tags, requestMessage)
}

// transport message is meant to carry original message using either encryption , signing or
func NewBinaryMessage(type_, service string, value []byte, props Props, tags Tags, requestMessage *FimpMessage) *FimpMessage {
	valEnc := base64.StdEncoding.EncodeToString(value)
	return NewMessage(type_, service, VTypeBinary, valEnc, props, tags, requestMessage)
}

func NewMessageFromBytes(msg []byte) (*FimpMessage, error) {
	fimpmsg := FimpMessage{}
	var err error
	fimpmsg.Type, err = jsonparser.GetString(msg, "type")
	fimpmsg.Service, err = jsonparser.GetString(msg, "serv")
	fimpmsg.ValueType, err = jsonparser.GetString(msg, "val_t")
	fimpmsg.UID, _ = jsonparser.GetString(msg, "uid")
	fimpmsg.CorrelationID, _ = jsonparser.GetString(msg, "corid")
	fimpmsg.CreationTime, _ = jsonparser.GetString(msg, "ctime")
	fimpmsg.ResponseToTopic, _ = jsonparser.GetString(msg, "resp_to")
	fimpmsg.Source, _ = jsonparser.GetString(msg, "src")

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
		var val []bool
		if _, err := jsonparser.ArrayEach(msg, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			item, _ := jsonparser.ParseBoolean(value)
			val = append(val, item)
		}, "val"); err != nil {
			return nil, err
		}

		fimpmsg.Value = val
	case VTypeStrArray:
		var val []string
		if _, err := jsonparser.ArrayEach(msg, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			item, _ := jsonparser.ParseString(value)
			val = append(val, item)
		}, "val"); err != nil {
			return nil, err
		}

		fimpmsg.Value = val
	case VTypeIntArray:
		var val []int64
		if _, err := jsonparser.ArrayEach(msg, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			item, _ := jsonparser.ParseInt(value)
			val = append(val, item)

		}, "val"); err != nil {
			return nil, err
		}
		fimpmsg.Value = val
	case VTypeFloatArray:
		var val []float64
		if _, err := jsonparser.ArrayEach(msg, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			item, _ := jsonparser.ParseFloat(value)
			val = append(val, item)
		}, "val"); err != nil {
			return nil, err
		}
		fimpmsg.Value = val

	case VTypeStrMap:
		val := make(map[string]string)
		if err := jsonparser.ObjectEach(msg, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			val[string(key)], err = jsonparser.ParseString(value)
			return nil
		}, "val"); err != nil {
			return nil, err
		}
		fimpmsg.Value = val

	case VTypeIntMap:
		val := make(map[string]int64)
		if err := jsonparser.ObjectEach(msg, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			val[string(key)], err = jsonparser.ParseInt(value)
			return nil
		}, "val"); err != nil {
			return nil, err
		}
		fimpmsg.Value = val

	case VTypeFloatMap:
		val := make(map[string]float64)
		if err := jsonparser.ObjectEach(msg, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			val[string(key)], err = jsonparser.ParseFloat(value)
			return nil
		}, "val"); err != nil {
			return nil, err
		}
		fimpmsg.Value = val

	case VTypeBoolMap:
		val := make(map[string]bool)
		if err := jsonparser.ObjectEach(msg, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			val[string(key)], err = jsonparser.ParseBoolean(value)
			return nil
		}, "val"); err != nil {
			return nil, err
		}
		fimpmsg.Value = val

	case VTypeBinary:
		fimpmsg.Value, err = jsonparser.GetString(msg, "val")
		//base64val, err := jsonparser.GetString(msg, "val")
		//if err != nil {
		//	return nil,err
		//}
		//fimpmsg.Value ,err = base64.StdEncoding.DecodeString(base64val)
		//if err != nil {
		//	return nil,err
		//}

	case VTypeObject:
		fimpmsg.ValueObj, _, _, err = jsonparser.Get(msg, "val")

	}
	if _, dt, _, err := jsonparser.Get(msg, "props"); dt != jsonparser.NotExist && dt != jsonparser.Null && err == nil {
		fimpmsg.Properties = make(Props)
		if err := jsonparser.ObjectEach(msg, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			fimpmsg.Properties[string(key)], err = jsonparser.ParseString(value)
			return nil
		}, "props"); err != nil {
			return nil, err
		}
	}

	return &fimpmsg, err

}
