package fimpgo

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/buger/jsonparser"
	"github.com/futurehomeno/fimpgo/fimptype"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

const (
	TimeFormat         = "2006-01-02T15:04:05.999Z07:00"
	VTypeString        = "string"
	VTypeInt           = "int"
	VTypeFloat         = "float"
	VTypeBool          = "bool"
	VTypeStrMap        = "str_map"
	VTypeIntMap        = "int_map"
	VTypeFloatMap      = "float_map"
	VTypeBoolMap       = "bool_map"
	VTypeStrArray      = "str_array"
	VTypeIntArray      = "int_array"
	VTypeFloatArray    = "float_array"
	VTypeBoolArray     = "bool_array"
	VTypeObject        = "object"
	VTypeBase64        = "base64"
	VTypeBinary        = "bin"
	VTypeNull          = "null"
	invalidValueFormat = "invalid value=%v type=%s exp=%T"

	Val = "val"
)

var timestampFormats = []string{
	time.RFC3339Nano,
	"2006-01-02T15:04:05.999999999Z0700",
	"2006-01-02 15:04:05.999999999 Z0700",
	"2006-01-02 15:04:05.999999999 Z07:00",
}

type Props map[string]string

func (p Props) GetIntValue(key string) (int, bool, error) {
	val, ok := p[key]
	if !ok {
		return 0, false, nil
	}

	i, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, true, fmt.Errorf("property %s value=%v invalid type exp=int got=%T", key, val, val)
	}

	return int(i), true, nil
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
		return 0, true, fmt.Errorf("property %s value=%v invalid type exp=float64 got=%T", key, val, val)
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
		return false, true, fmt.Errorf("property %s value=%v invalid type exp=bool got=%T", key, val, val)
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
		return time.Time{}, true, fmt.Errorf("property %s value=%v has invalid type exp=RFC3339 got=%T", key, val, val)
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
	StorageStrategyNone      StorageStrategy = "none"
)

type FimpMessage struct {
	Type            string                `json:"type"`
	Service         fimptype.ServiceNameT `json:"serv"`
	ValueType       string                `json:"val_t"`
	Value           any                   `json:"val"`
	ValueObj        []byte                `json:"-"`
	Tags            Tags                  `json:"tags"`
	Properties      Props                 `json:"props"`
	Storage         *Storage              `json:"storage,omitempty"`
	Version         string                `json:"ver"`
	CorrelationID   string                `json:"corid"`
	ResponseToTopic string                `json:"resp_to,omitempty"`
	Source          fimptype.ServiceNameT `json:"src,omitempty"`
	CreationTime    string                `json:"ctime"`
	UID             string                `json:"uid"`
	Topic           string                `json:"topic,omitempty"` // The field should be used to store original topic. It can be useful for converting message from MQTT to other transports.
}

func (msg *FimpMessage) SetValue(value any, valType string) {
	msg.Value = value
	msg.ValueType = valType
}

func (msg *FimpMessage) GetIntValue() (int, error) {
	switch val := msg.Value.(type) {
	case int64:
		if val > int64(math.MaxInt) || val < int64(math.MinInt) {
			return 0, fmt.Errorf("int64 value %d overflows int", val)
		}
		return int(val), nil
	case int:
		return val, nil
	case uint:
		if val > uint(math.MaxInt) {
			return 0, fmt.Errorf("uint value %d overflows int", val)
		}
		return int(val), nil
	}
	return 0, fmt.Errorf(invalidValueFormat, reflect.ValueOf(msg.Value), "int", msg.Value)
}

func (msg *FimpMessage) GetStringValue() (string, error) {
	val, ok := msg.Value.(string)
	if ok {
		return val, nil
	}
	return "", fmt.Errorf(invalidValueFormat, reflect.ValueOf(msg.Value), "string", msg.Value)
}

func (msg *FimpMessage) GetBoolValue() (bool, error) {
	val, ok := msg.Value.(bool)
	if ok {
		return val, nil
	}
	return false, fmt.Errorf(invalidValueFormat, reflect.ValueOf(msg.Value), "bool", msg.Value)
}

func (msg *FimpMessage) GetFloatValue() (float64, error) {
	val, ok := msg.Value.(float64)
	if ok {
		return val, nil
	}
	return 0, fmt.Errorf(invalidValueFormat, reflect.ValueOf(msg.Value), "float64", msg.Value)
}

func (msg *FimpMessage) GetStrArrayValue() ([]string, error) {
	val, ok := msg.Value.([]string)
	if ok {
		return val, nil
	}
	return nil, fmt.Errorf(invalidValueFormat, reflect.ValueOf(msg.Value), "[]string", msg.Value)
}

func (msg *FimpMessage) GetIntArrayValue() ([]int, error) {
	val64, ok := msg.Value.([]int64)
	if ok {
		ret := []int{}

		for _, v := range val64 {
			ret = append(ret, int(v))
		}
		return ret, nil
	}

	val, ok := msg.Value.([]int)
	if ok {
		return val, nil
	}

	return nil, fmt.Errorf(invalidValueFormat, reflect.ValueOf(msg.Value), "[]int", msg.Value)
}

func (msg *FimpMessage) GetFloatArrayValue() ([]float64, error) {
	val, ok := msg.Value.([]float64)
	if ok {
		return val, nil
	}
	return nil, fmt.Errorf(invalidValueFormat, reflect.ValueOf(msg.Value), "[]float64", msg.Value)
}

func (msg *FimpMessage) GetBoolArrayValue() ([]bool, error) {
	val, ok := msg.Value.([]bool)
	if ok {
		return val, nil
	}
	return nil, fmt.Errorf(invalidValueFormat, reflect.ValueOf(msg.Value), "[]bool", msg.Value)
}

func (msg *FimpMessage) GetStrMapValue() (map[string]string, error) {
	val, ok := msg.Value.(map[string]string)
	if ok {
		return val, nil
	}
	return nil, fmt.Errorf(invalidValueFormat, reflect.ValueOf(msg.Value), "map[string]string", msg.Value)
}

func (msg *FimpMessage) GetIntMapValue() (map[string]int, error) {
	val64, ok := msg.Value.(map[string]int64)
	if ok {
		ret := map[string]int{}

		for k, v := range val64 {
			ret[k] = int(v)
		}
		return ret, nil
	}

	val, ok := msg.Value.(map[string]int)
	if ok {
		return val, nil
	}
	return nil, fmt.Errorf(invalidValueFormat, reflect.ValueOf(msg.Value), "map[string]int", msg.Value)
}

func (msg *FimpMessage) GetFloatMapValue() (map[string]float64, error) {
	val, ok := msg.Value.(map[string]float64)
	if ok {
		return val, nil
	}
	return nil, fmt.Errorf(invalidValueFormat, reflect.ValueOf(msg.Value), "map[string]float64", msg.Value)
}

func (msg *FimpMessage) GetBoolMapValue() (map[string]bool, error) {
	val, ok := msg.Value.(map[string]bool)
	if ok {
		return val, nil
	}
	return nil, fmt.Errorf(invalidValueFormat, reflect.ValueOf(msg.Value), "map[string]bool", msg.Value)
}

func (msg *FimpMessage) GetRawObjectValue() []byte {
	return msg.ValueObj
}

func (msg *FimpMessage) GetObjectValue(objectBindVar any) error {
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

func (msg *FimpMessage) Str() string {
	var ret string

	if msg.Service != "" {
		if msg.Storage != nil && msg.Storage.Strategy != StorageStrategyNone {
			switch msg.Value.(type) {
			case float32, float64:
				ret = strings.TrimSpace(fmt.Sprintf("%s %s %s %.2f %s %s", msg.Source, msg.Service, msg.Type, msg.Value, msg.Storage.Strategy[:3], msg.Storage.SubValue))
			default:
				ret = strings.TrimSpace(fmt.Sprintf("%s %s %s %v %s %s", msg.Source, msg.Service, msg.Type, msg.Value, msg.Storage.Strategy[:3], msg.Storage.SubValue))
			}
		} else {
			switch msg.Value.(type) {
			case float32, float64:
				ret = fmt.Sprintf("%s %s %s %.2f", msg.Source, msg.Service, msg.Type, msg.Value)
			default:
				ret = fmt.Sprintf("%s %s %s %v", msg.Source, msg.Service, msg.Type, msg.Value)
			}
		}
	} else {
		ret = fmt.Sprintf("%s %s %v", msg.Source, msg.Type, msg.Value)
	}

	return ret
}

func NewMessage(type_ string, service fimptype.ServiceNameT, valueType string, value any, props Props, tags Tags, rqMsg *FimpMessage) *FimpMessage {
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

	if rqMsg != nil {
		msg.CorrelationID = rqMsg.UID
	}

	return &msg
}

func NewNullMessage(type_ string, service fimptype.ServiceNameT, props Props, tags Tags, rqMsg *FimpMessage) *FimpMessage {
	return NewMessage(type_, service, VTypeNull, nil, props, tags, rqMsg)
}

func NewStringMessage(type_ string, service fimptype.ServiceNameT, value string, props Props, tags Tags, rqMsg *FimpMessage) *FimpMessage {
	return NewMessage(type_, service, VTypeString, value, props, tags, rqMsg)
}

func NewIntMessage(type_ string, service fimptype.ServiceNameT, value int, props Props, tags Tags, rqMsg *FimpMessage) *FimpMessage {
	return NewMessage(type_, service, VTypeInt, value, props, tags, rqMsg)
}

func NewFloatMessage(type_ string, service fimptype.ServiceNameT, value float64, props Props, tags Tags, rqMsg *FimpMessage) *FimpMessage {
	return NewMessage(type_, service, VTypeFloat, value, props, tags, rqMsg)
}

func NewBoolMessage(type_ string, service fimptype.ServiceNameT, value bool, props Props, tags Tags, rqMsg *FimpMessage) *FimpMessage {
	return NewMessage(type_, service, VTypeBool, value, props, tags, rqMsg)
}

func NewStrArrayMessage(type_ string, service fimptype.ServiceNameT, value []string, props Props, tags Tags, rqMsg *FimpMessage) *FimpMessage {
	return NewMessage(type_, service, VTypeStrArray, value, props, tags, rqMsg)
}

func NewIntArrayMessage(type_ string, service fimptype.ServiceNameT, value []int, props Props, tags Tags, rqMsg *FimpMessage) *FimpMessage {
	return NewMessage(type_, service, VTypeIntArray, value, props, tags, rqMsg)
}

func NewFloatArrayMessage(type_ string, service fimptype.ServiceNameT, value []float64, props Props, tags Tags, rqMsg *FimpMessage) *FimpMessage {
	return NewMessage(type_, service, VTypeFloatArray, value, props, tags, rqMsg)
}

func NewBoolArrayMessage(type_ string, service fimptype.ServiceNameT, value []bool, props Props, tags Tags, rqMsg *FimpMessage) *FimpMessage {
	return NewMessage(type_, service, VTypeBoolArray, value, props, tags, rqMsg)
}

func NewStrMapMessage(type_ string, service fimptype.ServiceNameT, value map[string]string, props Props, tags Tags, rqMsg *FimpMessage) *FimpMessage {
	return NewMessage(type_, service, VTypeStrMap, value, props, tags, rqMsg)
}

func NewIntMapMessage(type_ string, service fimptype.ServiceNameT, value map[string]int, props Props, tags Tags, rqMsg *FimpMessage) *FimpMessage {
	return NewMessage(type_, service, VTypeIntMap, value, props, tags, rqMsg)
}

func NewFloatMapMessage(type_ string, service fimptype.ServiceNameT, value map[string]float64, props Props, tags Tags, rqMsg *FimpMessage) *FimpMessage {
	return NewMessage(type_, service, VTypeFloatMap, value, props, tags, rqMsg)
}

func NewBoolMapMessage(type_ string, service fimptype.ServiceNameT, value map[string]bool, props Props, tags Tags, rqMsg *FimpMessage) *FimpMessage {
	return NewMessage(type_, service, VTypeBoolMap, value, props, tags, rqMsg)
}

func NewObjectMessage(type_ string, service fimptype.ServiceNameT, value any, props Props, tags Tags, rqMsg *FimpMessage) *FimpMessage {
	return NewMessage(type_, service, VTypeObject, value, props, tags, rqMsg)
}

// NewBinaryMessage transport message is meant to carry original message using either encryption , signing or
func NewBinaryMessage(type_ string, service fimptype.ServiceNameT, value []byte, props Props, tags Tags, rqMsg *FimpMessage) *FimpMessage {
	valEnc := base64.StdEncoding.EncodeToString(value)
	return NewMessage(type_, service, VTypeBinary, valEnc, props, tags, rqMsg)
}

func NewMessageFromBytes(msg []byte) (*FimpMessage, error) { //nolint:gocyclo
	fimpmsg := FimpMessage{}
	var err error

	if fimpmsg.Type, err = jsonparser.GetString(msg, "type"); err != nil {
		log.Warnf("[fimpgo] NewMessageFromBytes type err: %v", err)
	}
	serviceStr, err := jsonparser.GetString(msg, "serv")

	if err != nil {
		log.Warnf("[fimpgo] NewMessageFromBytes serv err: %v", err)
	} else {
		fimpmsg.Service = fimptype.ServiceNameT(serviceStr)
	}
	if fimpmsg.ValueType, err = jsonparser.GetString(msg, "val_t"); err != nil {
		log.Warnf("[fimpgo] NewMessageFromBytes val_t err: %v", err)
	}
	if fimpmsg.UID, err = jsonparser.GetString(msg, "uid"); err != nil {
		log.Tracef("[fimpgo] NewMessageFromBytes uid err: %v", err)
	}
	if fimpmsg.CorrelationID, err = jsonparser.GetString(msg, "corid"); err != nil {
		log.Tracef("[fimpgo] NewMessageFromBytes coreid err: %v", err)
	}
	if fimpmsg.CreationTime, err = jsonparser.GetString(msg, "ctime"); err != nil {
		log.Tracef("[fimpgo] NewMessageFromBytes ctime err: %v", err)
	}
	if fimpmsg.ResponseToTopic, err = jsonparser.GetString(msg, "resp_to"); err != nil {
		log.Tracef("[fimpgo] NewMessageFromBytes resp_t err: %v", err)
	}

	sourceStr, err := jsonparser.GetString(msg, "src")

	if err != nil {
		log.Tracef("[fimpgo] NewMessageFromBytes src err: %v", err)
	} else {
		fimpmsg.Source = fimptype.ServiceNameT(sourceStr)
	}

	if fimpmsg.Topic, err = jsonparser.GetString(msg, "topic"); err != nil {
		log.Tracef("[fimpgo] NewMessageFromBytes topic err: %v", err)
	}
	if fimpmsg.Version, err = jsonparser.GetString(msg, "ver"); err != nil {
		log.Debugf("[fimpgo] NewMessageFromBytes ver err: %v", err)
	}

	err = nil

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
		_, err = jsonparser.ArrayEach(msg, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			item, e := jsonparser.ParseBoolean(value)
			if e != nil {
				log.Warnf("[fimpgo] Parse VTypeBoolArray err: %v", e)
				return
			}
			val = append(val, item)
		}, "val")

		fimpmsg.Value = val

	case VTypeStrArray:
		val := make([]string, 0)
		_, err = jsonparser.ArrayEach(msg, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			item, e := jsonparser.ParseString(value)
			if e != nil {
				log.Warnf("[fimpgo] Parse VTypeStrArray err: %v", e)
				return
			}
			val = append(val, item)
		}, "val")

		fimpmsg.Value = val

	case VTypeIntArray:
		val := make([]int, 0)
		_, err = jsonparser.ArrayEach(msg, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			item, e := jsonparser.ParseInt(value)
			if e != nil {
				log.Warnf("[fimpgo] Parse VTypeIntArray err: %v", e)
				return
			}
			val = append(val, int(item))
		}, "val")

		fimpmsg.Value = val

	case VTypeFloatArray:
		val := make([]float64, 0)
		_, err = jsonparser.ArrayEach(msg, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			item, e := jsonparser.ParseFloat(value)
			if e != nil {
				log.Warnf("[fimpgo] Parse VTypeFloatArray err: %v", e)
				return
			}
			val = append(val, item)
		}, "val")

		fimpmsg.Value = val

	case VTypeStrMap:
		val := make(map[string]string)
		err = jsonparser.ObjectEach(msg, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			tempStr, e := jsonparser.ParseString(value)
			if e != nil {
				log.Warnf("[fimpgo] Parse VTypeStrMap err: %v", e)
			} else {
				val[string(key)] = tempStr
			}
			return e
		}, "val")

		fimpmsg.Value = val

	case VTypeIntMap:
		val := make(map[string]int)
		err = jsonparser.ObjectEach(msg, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			tempInt, e := jsonparser.ParseInt(value)
			if e != nil {
				log.Warnf("[fimpgo] Parse VTypeIntMap err: %v", e)
			} else {
				val[string(key)] = int(tempInt)
			}
			return e
		}, "val")

		fimpmsg.Value = val

	case VTypeFloatMap:
		val := make(map[string]float64)
		err = jsonparser.ObjectEach(msg, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			tempFLoat, e := jsonparser.ParseFloat(value)
			if e != nil {
				log.Warnf("[fimpgo] Parse VTypeFloatMap err: %v", e)
			} else {
				val[string(key)] = tempFLoat
			}
			return e
		}, "val")

		fimpmsg.Value = val

	case VTypeBoolMap:
		val := make(map[string]bool)
		err = jsonparser.ObjectEach(msg, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			tempBool, e := jsonparser.ParseBoolean(value)
			if e != nil {
				log.Warnf("[fimpgo] Parse VTypeBoolMap err: %v", e)
			} else {
				val[string(key)] = tempBool
			}
			return e
		}, "val")

		fimpmsg.Value = val

	case VTypeBinary:
		fimpmsg.Value, err = jsonparser.GetString(msg, "val")
		if err != nil {
			log.Warnf("[fimpgo] GetString val err: %v", err)
		}

	case VTypeObject:
		fimpmsg.ValueObj, _, _, err = jsonparser.Get(msg, "val")
	case VTypeNull:
		fimpmsg.Value = nil
	default:
		return nil, jsonparser.UnknownValueTypeError
	}

	if err != nil {
		return nil, fmt.Errorf("val=%s err: %w", fimpmsg.ValueType, err)
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
