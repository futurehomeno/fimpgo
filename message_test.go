package fimpgo

import "testing"

func TestNewBoolMessage(t *testing.T) {
	msg := NewBoolMessage("cmd.binary.set", "out_bin_switch", true, nil, nil, nil)
	val, err := msg.GetBoolValue()
	if err != nil {
		t.Error(err)
	}
	if val == false {
		t.Error("Wrong value")
	}
	t.Log("ok")
}

func TestNewFloatMessage(t *testing.T) {

	msg := NewFloatMessage("evt.sensor.report", "temp_sensor", float64(35.5), nil, nil, nil)
	val, err := msg.GetFloatValue()
	if err != nil {
		t.Error(err)
	}
	if val != 35.5 {
		t.Error("Wrong value")
	}
	t.Log("ok")
}

func TestNewObjectMessage(t *testing.T) {

	type Event struct {
		Field1 int
		Field2 int
	}
	obj:= []Event{}

	obj = append(obj,Event{
		Field1: 1,
		Field2: 2,
	})
	msg := NewMessage("evt.timeline.report", "kind-owl",VTypeObject, obj, nil, nil, nil)
	bObj ,_ :=  msg.SerializeToJson()
	t.Log("ok",string(bObj))
}

func TestFimpMessage_SerializeBool(t *testing.T) {
	msg := NewBoolMessage("cmd.binary.set", "out_bin_switch", true, nil, nil, nil)
	serVal, err := msg.SerializeToJson()
	if err != nil {
		t.Error(err)
	}
	t.Log(string(serVal))
}

func TestFimpMessage_SerializeFloat(t *testing.T) {
	props := Props{}
	props["unit"] = "C"
	msg := NewFloatMessage("evt.sensor.report", "temp_sensor", float64(35.5), props, nil, nil)
	serVal, err := msg.SerializeToJson()
	if err != nil {
		t.Error(err)
	}
	t.Log(string(serVal))

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
	if err != nil {
		t.Log(err)
	}
	t.Log("ok")
}

func TestNewMessageFromBytes_BoolValue(t *testing.T) {
	msgString := "{\"serv\":\"out_bin_switch\",\"type\":\"cmd.binary.set\",\"val_t\":\"bool\",\"val\":true,\"props\":{\"p1\":\"pv1\"},\"tags\":null}"
	fimp, err := NewMessageFromBytes([]byte(msgString))
	if err != nil {
		t.Error(err)
	}
	val, err := fimp.GetBoolValue()
	if val != true {
		t.Error("Wrong value")
	}
	if fimp.Properties["p1"]!="pv1" {
		t.Error("Wrong props value")
	}
	t.Log("ok")
}

func TestNewMessageFromBytes_BoolInt(t *testing.T) {
	msgString := "{\"serv\":\"out_bin_switch\",\"type\":\"cmd.binary.set\",\"val_t\":\"int\",\"val\":1234,\"props\":null,\"tags\":null}"
	fimp, err := NewMessageFromBytes([]byte(msgString))
	if err != nil {
		t.Error(err)
	}
	val, err := fimp.GetIntValue()
	if val != 1234 {
		t.Error("Wrong value ",val)
	}
	t.Log("ok")
}

func TestFimpMessage_GetStrArrayValue(t *testing.T) {
	msgString := "{\"serv\":\"dev_sys\",\"type\":\"cmd.config.set\",\"val_t\":\"str_array\",\"val\":[\"val1\",\"val2\"],\"props\":null,\"tags\":null}"
	fimp, err := NewMessageFromBytes([]byte(msgString))
	if err != nil {
		t.Error(err)
	}

	val, err := fimp.GetStrArrayValue()
	if err != nil {
		t.Error(err)
	}
	if val[1] != "val2" {
		t.Error("Wrong map result : ",val[1])
	}
}

func TestFimpMessage_GetIntArrayValue(t *testing.T) {
	msgString := "{\"serv\":\"dev_sys\",\"type\":\"cmd.config.set\",\"val_t\":\"str_array\",\"val\":[123,1234],\"props\":null,\"tags\":null}"
	fimp, err := NewMessageFromBytes([]byte(msgString))
	if err != nil {
		t.Error(err)
	}

	val, err := fimp.GetIntArrayValue()
	if err != nil {
		t.Error(err)
	}
	if val[1] != 1234 {
		t.Error("Wrong map result : ",val[1])
	}
}

func TestFimpMessage_GetFloatArrayValue(t *testing.T) {
	msgString := "{\"serv\":\"dev_sys\",\"type\":\"cmd.config.set\",\"val_t\":\"float_array\",\"val\":[1.5,2.5],\"props\":null,\"tags\":null}"
	fimp, err := NewMessageFromBytes([]byte(msgString))
	if err != nil {
		t.Error(err)
	}

	val, err := fimp.GetFloatArrayValue()
	if err != nil {
		t.Error(err)
	}
	if val[1] != 2.5 {
		t.Error("Wrong map result : ",val[1])
	}
}

func TestFimpMessage_GetBoolArrayValue(t *testing.T) {
	msgString := "{\"serv\":\"dev_sys\",\"type\":\"cmd.config.set\",\"val_t\":\"bool_array\",\"val\":[true,true],\"props\":null,\"tags\":null}"
	fimp, err := NewMessageFromBytes([]byte(msgString))
	if err != nil {
		t.Error(err)
	}

	val, err := fimp.GetBoolArrayValue()
	if err != nil {
		t.Error(err)
	}
	if val[1] != true {
		t.Error("Wrong map result : ",val[1])
	}
}

func TestFimpMessage_GetStrMapValue(t *testing.T) {
	msgString := "{\"serv\":\"dev_sys\",\"type\":\"cmd.config.set\",\"val_t\":\"str_map\",\"val\":{\"param1\":\"val1\",\"param2\":\"val2\"},\"props\":null,\"tags\":null}"
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
	msgString := "{\"serv\":\"dev_sys\",\"type\":\"cmd.config.set\",\"val_t\":\"int_map\",\"val\":{\"param1\":1,\"param2\":2},\"props\":null,\"tags\":null}"
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
	msgString := "{\"serv\":\"dev_sys\",\"type\":\"cmd.config.set\",\"val_t\":\"float_map\",\"val\":{\"param1\":0.5,\"param2\":2.5},\"props\":null,\"tags\":null}"
	fimp, err := NewMessageFromBytes([]byte(msgString))
	if err != nil {
		t.Error(err)
	}

	val, err := fimp.GetFloatMapValue()
	if err != nil {
		t.Error(err)
	}
	if val["param2"] != 2.5 {
		t.Error("Wrong map result")
	}
}

func TestFimpMessage_GetBoolMapValue(t *testing.T) {
	msgString := "{\"serv\":\"dev_sys\",\"type\":\"cmd.config.set\",\"val_t\":\"bool_map\",\"val\":{\"param1\":true,\"param2\":true},\"props\":null,\"tags\":null}"
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

func BenchmarkFimpMessage_GetStrMapValue(b *testing.B) {
	msgString := []byte("{\"serv\":\"dev_sys\",\"type\":\"cmd.config.set\",\"val_t\":\"str_map\",\"val\":{\"param1\":\"val1\",\"param2\":\"val2\"},\"props\":null,\"tags\":null}")
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
	msgString := "{\"serv\":\"dev_sys\",\"type\":\"cmd.config.set\",\"val_t\":\"object\",\"val\":{\"param1\":\"val1\",\"param2\":\"val2\"},\"props\":null,\"tags\":null}"
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
}

func BenchmarkFimpMessage_GetObjectValue(b *testing.B) {
	type Config struct {
		Param1 string
		Param2 string
	}
	msgString := []byte("{\"serv\":\"dev_sys\",\"type\":\"cmd.config.set\",\"val_t\":\"object\",\"val\":{\"param1\":\"val1\",\"param2\":\"val2\"},\"props\":null,\"tags\":null}")
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
