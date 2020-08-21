package transport

import (
	"github.com/futurehomeno/fimpgo"
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
