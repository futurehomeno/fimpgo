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

func TestNewMessageFromBytes(t *testing.T) {
	msgString := "{\"serv\":\"out_bin_switch\",\"type\":\"cmd.binary.set\",\"val_t\":\"bool\",\"val\":true,\"props\":null,\"tags\":null}"
	fimp, err := NewMessageFromBytes([]byte(msgString))
	if err != nil {
		t.Error(err)
	}
	val, err := fimp.GetBoolValue()
	if val == false {
		t.Error("Wrong value")
	}
	t.Log("ok")
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
