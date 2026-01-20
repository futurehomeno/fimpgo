package fimpgo

import (
	"strings"
	"testing"
	"time"

	"github.com/buger/jsonparser"
	"github.com/stretchr/testify/assert"
)

func TestFimpMessage_SerializeToJson(t *testing.T) {
	tcs := []struct {
		name    string
		message *FimpMessage
	}{
		{
			name:    "Null message",
			message: NewNullMessage("test_type", "test_service", nil, nil, nil),
		},
		{
			name: "Null message with storage strategy",
			message: NewNullMessage("test_type", "test_service", nil, nil, nil).
				WithStorageStrategy(StorageStrategySkip, ""),
		},
		{
			name: "Null message with storage strategy, property and tag",
			message: NewNullMessage("test_type", "test_service", nil, nil, nil).
				WithProperty("prop1", "val1").
				WithTag("tag1").
				WithStorageStrategy(StorageStrategyAggregate, "val1"),
		},
	}

	for _, tc := range tcs {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			json, err := tc.message.SerializeToJson()
			assert.NoError(t, err)

			newMsg, err := NewMessageFromBytes(json)
			assert.NoError(t, err)

			assert.Equal(t, tc.message, newMsg)
		})
	}
}

func TestNewBoolMessage(t *testing.T) {
	msg := NewBoolMessage("cmd.binary.set", "out_bin_switch", true, nil, nil, nil)
	val, err := msg.GetBoolValue()
	if err != nil {
		t.Error(err)
	}
	if val == false {
		t.Error("Wrong value")
	}
}

func TestNewFloatMessage(t *testing.T) {

	msg := NewFloatMessage("evt.sensor.report", "temp_sensor", 35.5, nil, nil, nil)
	val, err := msg.GetFloatValue()

	if err != nil {
		t.Error(err)
	}
	if val != 35.5 {
		t.Error("Wrong value")
	}
}

func TestNewObjectMessage(t *testing.T) {
	type Event struct {
		Field1 int
		Field2 int
	}
	var obj []Event

	obj = append(obj, Event{
		Field1: 1,
		Field2: 2,
	})
	msg := NewMessage("evt.timeline.report", "kind-owl", VTypeObject, obj, nil, nil, nil)
	serVal, err := msg.SerializeToJson()

	if err != nil {
		t.Error(err)
	}

	if !strings.HasPrefix(string(serVal), `{"type":"evt.timeline.report","serv":"kind-owl","val_t":"object","val":[{"Field1":1,"Field2":2}],"tags":null,"props":null,"ver":"1","corid":"","ctime":"`) {
		t.Error("Serialization failed")
	}
}

func TestFimpMessage_SerializeBool(t *testing.T) {
	msg := NewBoolMessage("cmd.binary.set", "out_bin_switch", true, nil, nil, nil)
	serVal, err := msg.SerializeToJson()
	if err != nil {
		t.Error(err)
	}

	if !strings.HasPrefix(string(serVal), `{"type":"cmd.binary.set","serv":"out_bin_switch","val_t":"bool","val":true,"tags":null,"props":null,"ver":"1","corid":"","ctime":"`) {
		t.Error("Serialization failed")
	}
}

func TestFimpMessage_SerializeFloat(t *testing.T) {
	props := Props{}
	props["unit"] = "C"
	msg := NewFloatMessage("evt.sensor.report", "temp_sensor", 35.5, props, nil, nil)
	serVal, err := msg.SerializeToJson()
	if err != nil {
		t.Error(err)
	}

	if !strings.HasPrefix(string(serVal), `{"type":"evt.sensor.report","serv":"temp_sensor","val_t":"float","val":35.5,"tags":null,"props":{"unit":"C"},"ver":"1","corid":"","ctime":"`) {
		t.Error("Serialization failed")
	}
}

func BenchmarkFimpMessage_Serialize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		msg := NewBoolMessage("cmd.binary.set", "out_bin_switch", true, nil, nil, nil)
		_, err := msg.SerializeToJson()
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkFimpMessage_Serialize2(b *testing.B) {
	props := make(map[string]string)
	props["param1"] = "val1"
	for i := 0; i < b.N; i++ {
		msg := NewStrMapMessage("cmd.config.set", "dev_sys", props, nil, nil, nil)
		_, err := msg.SerializeToJson()
		if err != nil {
			b.Error(err)
		}
	}
}

func TestNewMessageFromBytes_CorruptedPayload1(t *testing.T) {
	msgString := "{123456789-=#$%"
	_, err := NewMessageFromBytes([]byte(msgString))
	assert.Equal(t, jsonparser.KeyPathNotFoundError, err)
}

func TestNewMessageFromBytes_BoolValue(t *testing.T) {
	msgString := `{"serv":"out_bin_switch","type":"cmd.binary.set","val_t":"bool","val":true,"props":{"p1":"pv1"},"tags":null}`
	fimp, err := NewMessageFromBytes([]byte(msgString))
	if err != nil {
		t.Error(err)
	}
	val, err := fimp.GetBoolValue()
	if val != true {
		t.Error("Wrong value")
	}
	if fimp.Properties["p1"] != "pv1" {
		t.Error("Wrong props value")
	}
}

func TestNewMessageFromBytes_BoolInt(t *testing.T) {
	msgString := `{"serv":"out_bin_switch","type":"cmd.binary.set","val_t":"int","val":1234,"props":null,"tags":null}`
	fimp, err := NewMessageFromBytes([]byte(msgString))
	if err != nil {
		t.Error(err)
	}
	val, err := fimp.GetIntValue()
	if val != 1234 {
		t.Error("Wrong value ", val)
	}
}

func TestNewMessageFromBytesWithProps(t *testing.T) {
	msgString := `{"serv":"out_bin_switch","type":"cmd.binary.set","val_t":"int","val":1234,"props":{"prop1":"val1"},"tags":null}`
	fimp, err := NewMessageFromBytes([]byte(msgString))
	if err != nil {
		t.Error(err)
	}
	val, err := fimp.GetIntValue()
	if val != 1234 {
		t.Error("Wrong value ", val)
	}
}

func TestFimpMessage_GetStrArrayValue(t *testing.T) {
	msgString := `{"serv":"dev_sys","type":"cmd.config.set","val_t":"str_array","val":["val1","val2"],"props":null,"tags":null}`
	fimp, err := NewMessageFromBytes([]byte(msgString))
	if err != nil {
		t.Error(err)
	}

	val, err := fimp.GetStrArrayValue()
	if err != nil {
		t.Error(err)
	}
	if val[1] != "val2" {
		t.Error("Wrong map result : ", val[1])
	}
}

func TestFimpMessage_GetIntArrayValue(t *testing.T) {
	msgString := `{"serv":"dev_sys","type":"cmd.config.set","val_t":"int_array","val":[123,1234],"props":null,"tags":null}`
	fimp, err := NewMessageFromBytes([]byte(msgString))
	if err != nil {
		t.Error(err)
	}

	val, err := fimp.GetIntArrayValue()
	if err != nil {
		t.Error(err)
	}
	if val[1] != 1234 {
		t.Error("Wrong map result : ", val[1])
	}
}

func TestFimpMessage_GetFloatArrayValue(t *testing.T) {
	msgString := `{"serv":"dev_sys","type":"cmd.config.set","val_t":"float_array","val":[1.5,2.5],"props":null,"tags":null}`
	fimp, err := NewMessageFromBytes([]byte(msgString))
	if err != nil {
		t.Error(err)
	}

	val, err := fimp.GetFloatArrayValue()
	if err != nil {
		t.Error(err)
	}
	if val[1] != 2.5 {
		t.Error("Wrong map result : ", val[1])
	}
}

func TestFimpMessage_GetBoolArrayValue(t *testing.T) {
	msgString := `{"serv":"dev_sys","type":"cmd.config.set","val_t":"bool_array","val":[true,true],"props":null,"tags":null}`
	fimp, err := NewMessageFromBytes([]byte(msgString))
	if err != nil {
		t.Error(err)
	}

	val, err := fimp.GetBoolArrayValue()
	if err != nil {
		t.Error(err)
	}
	if val[1] != true {
		t.Error("Wrong map result : ", val[1])
	}
}

func TestFimpMessage_GetStrMapValue(t *testing.T) {
	msgString := `{"serv":"dev_sys","type":"cmd.config.set","val_t":"str_map","val":{"param1":"val1","param2":"val2"},"props":null,"tags":null}`
	fimp, err := NewMessageFromBytes([]byte(msgString))
	if err != nil {
		t.Error(err)
	}

	val, err := fimp.GetStrMapValue()
	if err != nil {
		t.Error(err)
	}
	if val["param2"] != "val2" {
		t.Error("Wrong map result")
	}
}

func TestFimpMessage_GetIntMapValue(t *testing.T) {
	msgString := `{"serv":"dev_sys","type":"cmd.config.set","val_t":"int_map","val":{"param1":1,"param2":2},"props":null,"tags":null}`
	fimp, err := NewMessageFromBytes([]byte(msgString))
	if err != nil {
		t.Error(err)
	}

	val, err := fimp.GetIntMapValue()
	if err != nil {
		t.Error(err)
	}
	if val["param2"] != 2 {
		t.Error("Wrong map result")
	}
}

func TestFimpMessage_GetFloatMapValue(t *testing.T) {
	msgString := `{"serv":"dev_sys","type":"cmd.config.set","val_t":"float_map","val":{"param1":0.5,"param2":2.5,"param3":5},"props":null,"tags":null}`
	fimp, err := NewMessageFromBytes([]byte(msgString))
	if err != nil {
		t.Error(err)
	}

	val, err := fimp.GetFloatMapValue()
	if err != nil {
		t.Error(err)
	}
	if val["param2"] != 2.5 {
		t.Error("Wrong param2")
	}
	if val["param3"] != 5 {
		t.Error("Wrong param3")
	}
}

func TestFimpMessage_GetBoolMapValue(t *testing.T) {
	msgString := `{"serv":"dev_sys","type":"cmd.config.set","val_t":"bool_map","val":{"param1":true,"param2":true},"props":null,"tags":null}`
	fimp, err := NewMessageFromBytes([]byte(msgString))
	if err != nil {
		t.Error(err)
	}

	val, err := fimp.GetBoolMapValue()
	if err != nil {
		t.Error(err)
	}
	if val["param2"] != true {
		t.Error("Wrong map result")
	}
}

func TestProps_GetIntValue(t *testing.T) {
	msgString := `{"serv":"dev_sys","type":"cmd.config.set","val_t":"int","val":1234,"props":{"param1":"1","param2":"2"},"tags":null}`
	fimp, err := NewMessageFromBytes([]byte(msgString))
	if err != nil {
		t.Error(err)
	}

	props := fimp.Properties
	val, _, err := props.GetIntValue("param1")
	if err != nil {
		t.Error(err)
	}
	if val != 1 {
		t.Error("Wrong map result")
	}
}

func TestProps_GetStringValue(t *testing.T) {
	msgString := `{"serv":"dev_sys","type":"cmd.config.set","val_t":"str","val":"val1","props":{"param1":"val1","param2":"val2"},"tags":null}`
	fimp, err := NewMessageFromBytes([]byte(msgString))
	if err != nil {
		t.Error(err)
	}

	props := fimp.Properties
	val, _ := props.GetStringValue("param1")

	if val != "val1" {
		t.Error("Wrong map result")
	}
}

func TestProps_GetFloatValue(t *testing.T) {
	msgString := `{"serv":"dev_sys","type":"cmd.config.set","val_t":"float","val":1.5,"props":{"param1":"1.5","param2":"2.5"},"tags":null}`
	fimp, err := NewMessageFromBytes([]byte(msgString))
	if err != nil {
		t.Error(err)
	}

	props := fimp.Properties
	val, _, err := props.GetFloatValue("param1")
	if err != nil {
		t.Error(err)
	}
	if val != 1.5 {
		t.Error("Wrong map result")
	}
}

func TestProps_GetBoolValue(t *testing.T) {
	msgString := `{"serv":"dev_sys","type":"cmd.config.set","val_t":"bool","val":true,"props":{"param1":"true","param2":"false"},"tags":null}`
	fimp, err := NewMessageFromBytes([]byte(msgString))
	if err != nil {
		t.Error(err)
	}

	props := fimp.Properties
	val, _, err := props.GetBoolValue("param1")
	if err != nil {
		t.Error(err)
	}
	if val != true {
		t.Error("Wrong map result")
	}
}

func BenchmarkFimpMessage_GetStrMapValue(b *testing.B) {
	msgString := []byte(`{"serv":"dev_sys","type":"cmd.config.set","val_t":"str_map","val":{"param1":"val1","param2":"val2"},"props":null,"tags":null}`)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fimp, err := NewMessageFromBytes(msgString)
		if err != nil {
			b.Error(err)
		}

		val, err := fimp.GetStrMapValue()
		if err != nil {
			b.Error(err)
		}
		if val["param2"] != "val2" {
			b.Error("Wrong map result")
		}
	}
}

func TestFimpMessage_GetObjectValue(t *testing.T) {
	type Config struct {
		Param1 string
		Param2 string
	}
	msgString := `{"serv":"dev_sys","type":"cmd.config.set","val_t":"object","val":{"param1":"val1","param2":"val2"},"props":{"test":"1"},"tags":null}`
	fimp, err := NewMessageFromBytes([]byte(msgString))
	if err != nil {
		t.Error(err)
	}
	config := Config{}
	err = fimp.GetObjectValue(&config)
	if err != nil {
		t.Error(err)
	}
	if config.Param2 != "val2" {
		t.Error("Wrong map result")
	}
	binMsg, err := fimp.SerializeToJson()
	if err != nil {
		t.Error(err)
	}

	if !assert.Equal(t, string(binMsg), `{"type":"cmd.config.set","serv":"dev_sys","val_t":"object","val":{"param1":"val1","param2":"val2"},"tags":null,"props":{"test":"1"},"ver":"","corid":"","ctime":"","uid":""}`) {
		t.Error("Serialization failed")
	}
}

func BenchmarkFimpMessage_GetObjectValue(b *testing.B) {
	type Config struct {
		Param1 string
		Param2 string
	}
	msgString := []byte(`{"serv":"dev_sys","type":"cmd.config.set","val_t":"object","val":{"param1":"val1","param2":"val2"},"props":null,"tags":null}`)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fimp, err := NewMessageFromBytes(msgString)
		if err != nil {
			b.Error(err)
		}
		config := Config{}
		err = fimp.GetObjectValue(&config)
		if err != nil {
			b.Error(err)
		}
		if config.Param2 != "val2" {
			b.Error("Wrong map result")
		}
	}
}

func TestParseTime(t *testing.T) {
	t.Parallel()

	tt := []struct {
		timestamp string
		want      time.Time
	}{
		{timestamp: "2022-08-12T06:58:50.551867Z", want: time.Date(2022, 8, 12, 6, 58, 50, 551867000, time.UTC)},
		{timestamp: "2022-08-12T06:58:50Z", want: time.Date(2022, 8, 12, 6, 58, 50, 0, time.UTC)},
		{timestamp: "2022-08-12T08:58:53.383+02:00", want: time.Date(2022, 8, 12, 8, 58, 53, 383000000, time.FixedZone("", 7200))},
		{timestamp: "2022-08-12T08:58:53+02:00", want: time.Date(2022, 8, 12, 8, 58, 53, 0, time.FixedZone("", 7200))},
		{timestamp: "2022-08-12T08:58:53.551867+02:00", want: time.Date(2022, 8, 12, 8, 58, 53, 551867000, time.FixedZone("", 7200))},
		{timestamp: "2022-08-12T08:58:51+0200", want: time.Date(2022, 8, 12, 8, 58, 51, 0, time.FixedZone("", 7200))},
		{timestamp: "2022-08-12T08:58:51.383+0200", want: time.Date(2022, 8, 12, 8, 58, 51, 383000000, time.FixedZone("", 7200))},
		{timestamp: "2022-07-21 12:09:49 +0200", want: time.Date(2022, 7, 21, 12, 9, 49, 0, time.FixedZone("", 7200))},
		{timestamp: "2022-07-21 12:09:49.383 +0200", want: time.Date(2022, 7, 21, 12, 9, 49, 383000000, time.FixedZone("", 7200))},
		{timestamp: "2022-07-21 12:09:49 +02:00", want: time.Date(2022, 7, 21, 12, 9, 49, 0, time.FixedZone("", 7200))},
		{timestamp: "2022-07-21 12:09:49.383 +02:00", want: time.Date(2022, 7, 21, 12, 9, 49, 383000000, time.FixedZone("", 7200))},
		{timestamp: "invalid_format", want: time.Time{}},
		{timestamp: "", want: time.Time{}},
	}

	for _, tc := range tt {
		tc := tc

		t.Run(tc.timestamp, func(t *testing.T) {
			t.Parallel()

			got := ParseTime(tc.timestamp)

			assert.True(t, tc.want.Equal(got))
		})
	}
}
