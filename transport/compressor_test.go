package transport

import (
	"github.com/futurehomeno/fimpgo"
	"github.com/futurehomeno/fimpgo/fimptype"
	"testing"
)

func TestNewMsgCompressor(t *testing.T) {

	msg := fimpgo.NewFloatMessage("evt.sensor.report", "temp_sensor", 35.5, nil, nil, nil)

	comp :=  NewMsgCompressor("","")
	decomp := NewMsgCompressor("","")

	for i:=0;i<5;i++ {
		bmsg,_:=msg.SerializeToJson()

		t.Log("Uncompressed message size = ",len(bmsg))
		compMsg,err := comp.CompressFimpMsg(msg)
		if err != nil {
			t.Fatal("Compressor error :",err)
		}
		t.Log("Compressed message size = ",len(compMsg))
		fimpMsg , err := decomp.DecompressFimpMsg(compMsg)
		if err != nil {
			t.Fatal("Compressor error :",err)
		}
		if fimpMsg.Service != "temp_sensor" {
			t.Fatal("Incorrect service name ")
		}else {
			t.Log("All good")
		}
	}
}

func TestNewMsgCompressor2(t *testing.T) {
	val := fimptype.ThingInclusionReport{}
	val.Alias= "test_alias"
	msg0 := fimpgo.NewObjectMessage("evt.inclusion.report", "test", val, nil, nil, nil)
	bmsg0,_ := msg0.SerializeToJson()



	comp :=  NewMsgCompressor("","")
	decomp := NewMsgCompressor("","")

	for i:=0;i<10;i++ {
		msg,err := fimpgo.NewMessageFromBytes(bmsg0)
		if err != nil {
			t.Fatal("Desirialization error")
		}
		msg.Topic = "some/topic"
		if msg.ValueType == fimpgo.VTypeObject {
			err := msg.GetObjectValue(&msg.Value)
			if err != nil {
				t.Fatal("<ses> Compression fimp error:",err.Error())
			}
		}
		bmsg,_:=msg.SerializeToJson()
		t.Log("Uncompressed message size = ",len(bmsg))
		compMsg,err := comp.CompressFimpMsg(msg)
		if err != nil {
			t.Fatal("Compressor error :",err)
		}
		t.Log("Compressed message size = ",len(compMsg))
		fimpMsg , err := decomp.DecompressFimpMsg(compMsg)
		if err != nil {
			t.Fatal("Compressor error :",err)
		}
		v1 := fimptype.ThingInclusionReport{}
		fimpMsg.GetObjectValue(&v1)
		if fimpMsg.Service != "test" || v1.Alias != "test_alias" {
			t.Fatal("Incorrect service name ")
		}else {
			t.Log("All good")
		}
	}
}

