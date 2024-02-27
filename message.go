package fimpgo

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/buger/jsonparser"
	"github.com/google/uuid"
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

var timestampFormats = []string{
	time.RFC3339Nano,
	"2006-01-02T15:04:05.999999999Z0700",
	"2006-01-02 15:04:05.999999999 Z0700",
	"2006-01-02 15:04:05.999999999 Z07:00",
}

type Props map[string]string

func (p Props) GetIntValue(key string) (int64, bool, error) {
	val, ok := p[key]
	if !ok {
		return 0, false, nil
	}

	i, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, true, fmt.Errorf("property %s has wrong value type, expected int, got %s", key, val)
	}

	return i, true, nil
}

func (p Props) GetStringValue(key string) (string, bool) {
	val, ok := p[key]
	if !ok {
		return "", false
	}

	return val, true
}

func (p Props) GetFloatValue(key string) (float64, bool, error) {
	val, ok := p[key]
	if !ok {
		return 0, false, nil
	}

	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0, true, fmt.Errorf("property %s has wrong value type, expected float, got %s", key, val)
	}

	return f, true, nil
}

func (p Props) GetBoolValue(key string) (bool, bool, error) {
	val, ok := p[key]
	if !ok {
		return false, false, nil
	}

	b, err := strconv.ParseBool(val)
	if err != nil {
		return false, true, fmt.Errorf("property %s has wrong value type, expected bool, got %s", key, val)
	}

	return b, true, nil
}

func (p Props) GetTimestampValue(key string) (time.Time, bool, error) {
	val, ok := p[key]
	if !ok {
		return time.Time{}, false, nil
	}

	t := ParseTime(val)
	if t.IsZero() {
		return time.Time{}, true, fmt.Errorf("property %s has wrong value type, expected RFC3339 timestamp, got %s", key, val)
	}

	return t, true, nil
}

type Tags []string

// Storage is used to define optional message storage strategy.
type Storage struct {
	Strategy StorageStrategy `json:"strategy,omitempty"`
	SubValue string          `json:"sub_value,omitempty"`
}

// StorageStrategy defines message storage strategy.
type StorageStrategy string

// Constants defining storage strategies.
const (
	StorageStrategyAggregate StorageStrategy = "aggregate"
	StorageStrategySkip      StorageStrategy = "skip"
	StorageStrategySplit     StorageStrategy = "split"
)

type FimpMessage struct {
	Type            string      `json:"type"`
	Service         string      `json:"serv"`
	ValueType       string      `json:"val_t"`
	Value           interface{} `json:"val"`
	ValueObj        []byte      `json:"-"`
	Tags            Tags        `json:"tags"`
	Properties      Props       `json:"props"`
	Storage         *Storage    `json:"storage,omitempty"`
	Version         string      `json:"ver"`
	CorrelationID   string      `json:"corid"`
	ResponseToTopic string      `json:"resp_to,omitempty"`
	Source          string      `json:"src,omitempty"`
	CreationTime    string      `json:"ctime"`
	UID             string      `json:"uid"`
	Topic           string      `json:"topic,omitempty"` // The field should be used to store original topic. It can be useful for converting message from MQTT to other transports.
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
	if msg.ValueType == VTypeObject {
		if msg.Value == nil && msg.ValueObj != nil {
			// This is for object pass though.
			jsonBA, err = jsonparser.Set(jsonBA, msg.ValueObj, "val")
		}
	}
	return jsonBA, err

}

// GetCreationTime returns parsed creation time of the message.
func (msg *FimpMessage) GetCreationTime() time.Time {
	return ParseTime(msg.CreationTime)
}

// WithStorageStrategy sets storage strategy for the message.
func (msg *FimpMessage) WithStorageStrategy(strategy StorageStrategy, subValue string) *FimpMessage {
	msg.Storage = &Storage{Strategy: strategy, SubValue: subValue}

	return msg
}

// WithProperty sets property for the message.
func (msg *FimpMessage) WithProperty(property, value string) *FimpMessage {
	if msg.Properties == nil {
		msg.Properties = make(Props)
	}

	msg.Properties[property] = value

	return msg
}

// WithTag adds tag to the message.
func (msg *FimpMessage) WithTag(tag string) *FimpMessage {
	msg.Tags = append(msg.Tags, tag)

	return msg
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
		UID:          uuid.New().String(),
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

// NewBinaryMessage transport message is meant to carry original message using either encryption , signing or
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
	fimpmsg.Topic, _ = jsonparser.GetString(msg, "topic")
	fimpmsg.Version, _ = jsonparser.GetString(msg, "ver")

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
		val := make([]bool, 0)
		if _, err := jsonparser.ArrayEach(msg, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			item, _ := jsonparser.ParseBoolean(value)
			val = append(val, item)
		}, "val"); err != nil {
			return nil, err
		}

		fimpmsg.Value = val
	case VTypeStrArray:
		val := make([]string, 0)
		if _, err := jsonparser.ArrayEach(msg, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			item, _ := jsonparser.ParseString(value)
			val = append(val, item)
		}, "val"); err != nil {
			return nil, err
		}

		fimpmsg.Value = val
	case VTypeIntArray:
		val := make([]int64, 0)
		if _, err := jsonparser.ArrayEach(msg, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			item, _ := jsonparser.ParseInt(value)
			val = append(val, item)

		}, "val"); err != nil {
			return nil, err
		}
		fimpmsg.Value = val
	case VTypeFloatArray:
		val := make([]float64, 0)
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

	if properties, dt, _, err := jsonparser.Get(msg, "props"); dt != jsonparser.NotExist && dt != jsonparser.Null && err == nil {
		err := json.Unmarshal(properties, &fimpmsg.Properties)
		if err != nil {
			return nil, err
		}
	}

	if storage, dt, _, err := jsonparser.Get(msg, "storage"); dt != jsonparser.NotExist && dt != jsonparser.Null && err == nil {
		err := json.Unmarshal(storage, &fimpmsg.Storage)
		if err != nil {
			return nil, err
		}
	}

	if tags, dt, _, err := jsonparser.Get(msg, "tags"); dt != jsonparser.NotExist && dt != jsonparser.Null && err == nil {
		err := json.Unmarshal(tags, &fimpmsg.Tags)
		if err != nil {
			return nil, err
		}
	}

	return &fimpmsg, err
}

// ParseTime is a helper function to parse a timestamp from a string from various variations of RFC3339.
func ParseTime(timestamp string) time.Time {
	for _, format := range timestampFormats {
		t, err := time.Parse(format, timestamp)
		if err == nil {
			return t
		}
	}

	return time.Time{}
}
