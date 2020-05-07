package transport

import (
	"github.com/futurehomeno/fimpgo"
	"testing"
)

func TestSignMessageES256(t *testing.T) {
	keys := NewEcKeyPair()
	if err:=keys.Generate();err !=nil {
		t.Error("Key generation error",err)
		t.FailNow()
	}
	private , pub := keys.GetEncodedKeys()
	t.Log(private)
	t.Log(pub)
	msg := fimpgo.NewFloatMessage("evt.sensor.report", "temp_sensor", 35.5, nil, nil, nil)
	got, err := SignMessageES256(msg, nil, "alex@gmail.com", keys, nil)
	if err = keys.Generate();err !=nil {
		t.Error("Key generation error",err)
		t.FailNow()
	}

	bmsg,err := got.SerializeToJson()
	if err = keys.Generate();err !=nil {
		t.Error("Key generation error",err)
		t.FailNow()
	}
	t.Log(string(bmsg))

	keys2 := NewEcKeyPair()
	err = keys2.ImportPublicKey(pub)
	if err = keys.Generate();err !=nil {
		t.Error("Failed to import the key",err)
		t.FailNow()
	}
	result,err := GetVerifiedMessageES256(got,keys2)
	if err != nil || result == nil {
		t.Error("Signature is not valid")
	}else {
		t.Log(result)
	}

}