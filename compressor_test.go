package fimpgo

import (
	"testing"

	"github.com/futurehomeno/fimpgo/fimptype"
)

func TestNewMsgCompressor(t *testing.T) {
	msg := NewFloatMessage("evt.sensor.report", "temp_sensor", 35.5, nil, nil, nil)

	comp := NewMsgCompressor("", "")
	decomp := NewMsgCompressor("", "")

	for i := 0; i < 5; i++ {
		compMsg, err := comp.CompressFimpMsg(msg)
		if err != nil {
			t.Fatal("Compressor error :", err)
		}

		fimpMsg, err := decomp.DecompressFimpMsg(compMsg)
		if err != nil {
			t.Fatal("Compressor error 2 :", err)
		}
		if fimpMsg.Service != "temp_sensor" {
			t.Fatal("Incorrect service name ", fimpMsg.Service)
		}
	}
}

func TestNewMsgCompressor2(t *testing.T) {
	val := fimptype.ThingInclusionReport{}
	val.HwVersion = "hw_version"
	msg0 := NewObjectMessage("evt.inclusion.report", "test", val, nil, nil, nil)
	bmsg0, _ := msg0.SerializeToJson()

	comp := NewMsgCompressor("", "")
	decomp := NewMsgCompressor("", "")

	for range 10 {
		msg, err := NewMessageFromBytes(bmsg0)
		if err != nil {
			t.Fatal("Deserialization error")
		}
		msg.Topic = "some/topic"
		if msg.ValueType == VTypeObject {
			err := msg.GetObjectValue(&msg.Value)
			if err != nil {
				t.Fatal("Compression fimp error:", err.Error())
			}
		}
		bmsg, _ := msg.SerializeToJson()

		compMsg, err := comp.CompressBinMsg(bmsg)
		if err != nil {
			t.Fatal("Compressor error :", err)
		}

		fimpMsgBin, err := decomp.DecompressBinMsg(compMsg)
		if err != nil {
			t.Fatal("Compressor error 1:", err)
		}

		fimpMsg, err := NewMessageFromBytes(fimpMsgBin)
		//fimpMsg , err := decomp.DecompressFimpMsg(compMsg)
		if err != nil {
			t.Fatal("Compressor error 2:", err)
		}
		v1 := fimptype.ThingInclusionReport{}
		if err := fimpMsg.GetObjectValue(&v1); err != nil {
			t.Fatal("Getting object value error:", err)
		}

		if fimpMsg.Service != "test" || v1.HwVersion != "hw_version" {
			t.Fatalf("Incorrect service name %s or hw version %s", fimpMsg.Service, v1.HwVersion)
		}
	}
}
