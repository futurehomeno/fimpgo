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
	invalidValueFormat = "invalid value=%v type=%s exp=%T"
	ValField           = "val"
	ValTypeField       = "val_t"
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
	Interface       string                 `json:"type"`
	Service         fimptype.ServiceNameT  `json:"serv"`
	ValueType       fimptype.ValueTypeT    `json:"val_t"`
	Value           any                    `json:"val"`
	ValueObj        []byte                 `json:"-"`
	Tags            Tags                   `json:"tags"`
	Properties      Props                  `json:"props"`
	Storage         *Storage               `json:"storage,omitempty"`
	Version         string                 `json:"ver"`
	CorrelationID   string                 `json:"corid"`
	ResponseToTopic string                 `json:"resp_to,omitempty"`
	Source          fimptype.ResourceNameT `json:"src,omitempty"`
	CreationTime    string                 `json:"ctime"`
	UID             string                 `json:"uid"`
	Topic           string                 `json:"topic,omitempty"` // The field should be used to store original topic. It can be useful for converting message from MQTT to other transports.
}

func (msg *FimpMessage) SetValue(value any, valType fimptype.ValueTypeT) {
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
	if msg.ValueType == fimptype.VTypeObject {
		if msg.Value == nil && msg.ValueObj != nil {
			// This is for object pass though.
			jsonBA, err = jsonparser.Set(jsonBA, msg.ValueObj, ValField)
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
			strategyStr := string(msg.Storage.Strategy)
			if len(strategyStr) >= 3 {
				strategyStr = strategyStr[:3]
			}

			switch msg.Value.(type) {
			case float32, float64:
				ret = strings.TrimSpace(fmt.Sprintf("%s %s %s %.2f %s %s", msg.Source, msg.Service, msg.Interface, msg.Value, strategyStr, msg.Storage.SubValue))
			default:
				ret = strings.TrimSpace(fmt.Sprintf("%s %s %s %v %s %s", msg.Source, msg.Service, msg.Interface, msg.Value, strategyStr, msg.Storage.SubValue))
			}
		} else {
			switch msg.Value.(type) {
			case float32, float64:
				ret = fmt.Sprintf("%s %s %s %.2f", msg.Source, msg.Service, msg.Interface, msg.Value)
			default:
				ret = fmt.Sprintf("%s %s %s %v", msg.Source, msg.Service, msg.Interface, msg.Value)
			}
		}
	} else {
		ret = fmt.Sprintf("%s %s %v", msg.Source, msg.Interface, msg.Value)
	}

	return ret
}

func NewMessage(iface string, service fimptype.ServiceNameT, valueType fimptype.ValueTypeT, value any, props Props, tags Tags, rqMsg *FimpMessage) *FimpMessage {
	msg := FimpMessage{
		Interface:    iface,
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

func NewNullMessage(iface string, service fimptype.ServiceNameT, props Props, tags Tags, rqMsg *FimpMessage) *FimpMessage {
	return NewMessage(iface, service, fimptype.VTypeNull, nil, props, tags, rqMsg)
}

func NewStringMessage(iface string, service fimptype.ServiceNameT, value string, props Props, tags Tags, rqMsg *FimpMessage) *FimpMessage {
	return NewMessage(iface, service, fimptype.VTypeString, value, props, tags, rqMsg)
}

func NewIntMessage(iface string, service fimptype.ServiceNameT, value int, props Props, tags Tags, rqMsg *FimpMessage) *FimpMessage {
	return NewMessage(iface, service, fimptype.VTypeInt, value, props, tags, rqMsg)
}

func NewFloatMessage(iface string, service fimptype.ServiceNameT, value float64, props Props, tags Tags, rqMsg *FimpMessage) *FimpMessage {
	return NewMessage(iface, service, fimptype.VTypeFloat, value, props, tags, rqMsg)
}

func NewBoolMessage(iface string, service fimptype.ServiceNameT, value bool, props Props, tags Tags, rqMsg *FimpMessage) *FimpMessage {
	return NewMessage(iface, service, fimptype.VTypeBool, value, props, tags, rqMsg)
}

func NewStrArrayMessage(iface string, service fimptype.ServiceNameT, value []string, props Props, tags Tags, rqMsg *FimpMessage) *FimpMessage {
	return NewMessage(iface, service, fimptype.VTypeStrArray, value, props, tags, rqMsg)
}

func NewIntArrayMessage(iface string, service fimptype.ServiceNameT, value []int, props Props, tags Tags, rqMsg *FimpMessage) *FimpMessage {
	return NewMessage(iface, service, fimptype.VTypeIntArray, value, props, tags, rqMsg)
}

func NewFloatArrayMessage(iface string, service fimptype.ServiceNameT, value []float64, props Props, tags Tags, rqMsg *FimpMessage) *FimpMessage {
	return NewMessage(iface, service, fimptype.VTypeFloatArray, value, props, tags, rqMsg)
}

func NewBoolArrayMessage(iface string, service fimptype.ServiceNameT, value []bool, props Props, tags Tags, rqMsg *FimpMessage) *FimpMessage {
	return NewMessage(iface, service, fimptype.VTypeBoolArray, value, props, tags, rqMsg)
}

func NewStrMapMessage(iface string, service fimptype.ServiceNameT, value map[string]string, props Props, tags Tags, rqMsg *FimpMessage) *FimpMessage {
	return NewMessage(iface, service, fimptype.VTypeStrMap, value, props, tags, rqMsg)
}

func NewIntMapMessage(iface string, service fimptype.ServiceNameT, value map[string]int, props Props, tags Tags, rqMsg *FimpMessage) *FimpMessage {
	return NewMessage(iface, service, fimptype.VTypeIntMap, value, props, tags, rqMsg)
}

func NewFloatMapMessage(iface string, service fimptype.ServiceNameT, value map[string]float64, props Props, tags Tags, rqMsg *FimpMessage) *FimpMessage {
	return NewMessage(iface, service, fimptype.VTypeFloatMap, value, props, tags, rqMsg)
}

func NewBoolMapMessage(iface string, service fimptype.ServiceNameT, value map[string]bool, props Props, tags Tags, rqMsg *FimpMessage) *FimpMessage {
	return NewMessage(iface, service, fimptype.VTypeBoolMap, value, props, tags, rqMsg)
}

func NewObjectMessage(iface string, service fimptype.ServiceNameT, value any, props Props, tags Tags, rqMsg *FimpMessage) *FimpMessage {
	return NewMessage(iface, service, fimptype.VTypeObject, value, props, tags, rqMsg)
}

// NewBinaryMessage transport message is meant to carry original message using either encryption , signing or
func NewBinaryMessage(iface string, service fimptype.ServiceNameT, value []byte, props Props, tags Tags, rqMsg *FimpMessage) *FimpMessage {
	valEnc := base64.StdEncoding.EncodeToString(value)
	return NewMessage(iface, service, fimptype.VTypeBinary, valEnc, props, tags, rqMsg)
}

func NewMessageFromBytes(msg []byte) (*FimpMessage, error) { //nolint:gocyclo
	fimpmsg := FimpMessage{}
	var err error

	if fimpmsg.Interface, err = jsonparser.GetString(msg, "type"); err != nil {
		log.Warnf("[fimpgo] Parse %s type err: %v", fimpmsg.Interface, err)
		return nil, err
	}

	serviceStr, err := jsonparser.GetString(msg, "serv")

	if err != nil {
		log.Warnf("[fimpgo] Parse %s serv err: %v", fimpmsg.Interface, err)
	} else {
		fimpmsg.Service = fimptype.ServiceNameT(serviceStr)
	}

	valueTypeStr, err := jsonparser.GetString(msg, ValTypeField)

	if err != nil {
		log.Warnf("[fimpgo] Parse %s %s err: %v", fimpmsg.Interface, ValTypeField, err)
	} else {
		fimpmsg.ValueType = fimptype.ValueTypeT(valueTypeStr)
	}

	if fimpmsg.UID, err = jsonparser.GetString(msg, "uid"); err != nil {
		log.Tracef("[fimpgo] Parse %s uid err: %v", fimpmsg.Interface, err)
	}
	if fimpmsg.CorrelationID, err = jsonparser.GetString(msg, "corid"); err != nil {
		log.Tracef("[fimpgo] Parse %s coreid err: %v", fimpmsg.Interface, err)
	}
	if fimpmsg.CreationTime, err = jsonparser.GetString(msg, "ctime"); err != nil {
		log.Tracef("[fimpgo] Parse %s ctime err: %v", fimpmsg.Interface, err)
	}
	if fimpmsg.ResponseToTopic, err = jsonparser.GetString(msg, "resp_to"); err != nil {
		log.Tracef("[fimpgo] Parse %s resp_t err: %v", fimpmsg.Interface, err)
	}

	sourceStr, err := jsonparser.GetString(msg, "src")

	if err != nil {
		log.Tracef("[fimpgo] Parse %s src err: %v", fimpmsg.Interface, err)
	} else {
		fimpmsg.Source = fimptype.ResourceNameT(sourceStr)
	}

	if fimpmsg.Topic, err = jsonparser.GetString(msg, "topic"); err != nil {
		log.Tracef("[fimpgo] Parse %s topic err: %v", fimpmsg.Interface, err)
	}

	if fimpmsg.Version, err = jsonparser.GetString(msg, "ver"); err != nil {
		log.Tracef("[fimpgo] Parse %s ver err: %v", fimpmsg.Interface, err)
	}

	err = nil

	switch fimpmsg.ValueType {
	case fimptype.VTypeString:
		fimpmsg.Value, err = jsonparser.GetString(msg, ValField)
	case fimptype.VTypeBool:
		fimpmsg.Value, err = jsonparser.GetBoolean(msg, ValField)
	case fimptype.VTypeInt:
		fimpmsg.Value, err = jsonparser.GetInt(msg, ValField)
	case fimptype.VTypeFloat:
		fimpmsg.Value, err = jsonparser.GetFloat(msg, ValField)
	case fimptype.VTypeBoolArray:
		val := make([]bool, 0)
		_, err = jsonparser.ArrayEach(msg, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			item, e := jsonparser.ParseBoolean(value)
			if e != nil {
				log.Warnf("[fimpgo] Parse %s VTypeBoolArray err: %v", fimpmsg.Interface, e)
				return
			}
			val = append(val, item)
		}, ValField)

		fimpmsg.Value = val

	case fimptype.VTypeStrArray:
		val := make([]string, 0)
		_, err = jsonparser.ArrayEach(msg, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			item, e := jsonparser.ParseString(value)
			if e != nil {
				log.Warnf("[fimpgo] Parse %s VTypeStrArray err: %v", fimpmsg.Interface, e)
				return
			}
			val = append(val, item)
		}, ValField)

		fimpmsg.Value = val

	case fimptype.VTypeIntArray:
		val := make([]int, 0)
		_, err = jsonparser.ArrayEach(msg, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			item, e := jsonparser.ParseInt(value)
			if e != nil {
				log.Warnf("[fimpgo] Parse %s VTypeIntArray err: %v", fimpmsg.Interface, e)
				return
			}
			val = append(val, int(item))
		}, ValField)

		fimpmsg.Value = val

	case fimptype.VTypeFloatArray:
		val := make([]float64, 0)
		_, err = jsonparser.ArrayEach(msg, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			item, e := jsonparser.ParseFloat(value)
			if e != nil {
				log.Warnf("[fimpgo] Parse %s VTypeFloatArray err: %v", fimpmsg.Interface, e)
				return
			}
			val = append(val, item)
		}, ValField)

		fimpmsg.Value = val

	case fimptype.VTypeStrMap:
		val := make(map[string]string)
		err = jsonparser.ObjectEach(msg, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			tempStr, e := jsonparser.ParseString(value)
			if e != nil {
				log.Warnf("[fimpgo] Parse %s VTypeStrMap err: %v", fimpmsg.Interface, e)
			} else {
				val[string(key)] = tempStr
			}
			return e
		}, ValField)

		fimpmsg.Value = val

	case fimptype.VTypeIntMap:
		val := make(map[string]int)
		err = jsonparser.ObjectEach(msg, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			tempInt, e := jsonparser.ParseInt(value)
			if e != nil {
				log.Warnf("[fimpgo] Parse %s VTypeIntMap err: %v", fimpmsg.Interface, e)
			} else {
				val[string(key)] = int(tempInt)
			}
			return e
		}, ValField)

		fimpmsg.Value = val

	case fimptype.VTypeFloatMap:
		val := make(map[string]float64)
		err = jsonparser.ObjectEach(msg, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			tempFLoat, e := jsonparser.ParseFloat(value)
			if e != nil {
				log.Warnf("[fimpgo] Parse %s VTypeFloatMap err: %v", fimpmsg.Interface, e)
			} else {
				val[string(key)] = tempFLoat
			}
			return e
		}, ValField)

		fimpmsg.Value = val

	case fimptype.VTypeBoolMap:
		val := make(map[string]bool)
		err = jsonparser.ObjectEach(msg, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			tempBool, e := jsonparser.ParseBoolean(value)
			if e != nil {
				log.Warnf("[fimpgo] Parse %s VTypeBoolMap err: %v", fimpmsg.Interface, e)
			} else {
				val[string(key)] = tempBool
			}
			return e
		}, ValField)

		fimpmsg.Value = val

	case fimptype.VTypeBinary:
		fimpmsg.Value, err = jsonparser.GetString(msg, ValField)
		if err != nil {
			log.Warnf("[fimpgo] GetString %s val err: %v", fimpmsg.Interface, err)
		}

	case fimptype.VTypeObject:
		fimpmsg.ValueObj, _, _, err = jsonparser.Get(msg, ValField)
	case fimptype.VTypeNull:
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
