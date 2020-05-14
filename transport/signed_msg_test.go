package transport

import (
	"encoding/json"
	"github.com/futurehomeno/fimpgo"
	"github.com/futurehomeno/fimpgo/security"
	"testing"
)

func TestSignMessageES256(t *testing.T) {
	keys := integration.NewEcdsaKey()
	if err:=keys.Generate();err !=nil {
		t.Error("SerializedKey generation error",err)
		t.FailNow()
	}
	private , pub := keys.ExportX509EncodedKeys()
	t.Log(private)
	t.Log(pub)
	msg := fimpgo.NewFloatMessage("evt.sensor.report", "temp_sensor", 35.5, nil, nil, nil)
	got, err := SignMessageES256(msg, nil, "alex@gmail.com", keys, nil)
	if err !=nil {
		t.Error("Signing error",err)
		t.FailNow()
	}

	bmsg,err := got.SerializeToJson()
	if err = keys.Generate();err !=nil {
		t.Error("SerializedKey generation error",err)
		t.FailNow()
	}
	t.Log("X:",keys.PublicKey().X.Text(16))
	t.Log("Y:",keys.PublicKey().Y.Text(16))

	t.Log(string(bmsg))

	keys2 := integration.NewEcdsaKey()
	err = keys2.ImportX509PublicKey(pub)
	if err != nil {
		t.Error("Wrong key")
		t.FailNow()
	}

	result,err := GetVerifiedMessageES256(got,keys2)
	if err != nil || result == nil {
		t.Error("Signature is not valid")
	}else {
		t.Log(result)
	}

}

func TestSignMessageES256_TestKey(t *testing.T) {
	keys := integration.NewEcdsaKey()
	if err:=keys.Generate();err !=nil {
		t.Error("SerializedKey generation error",err)
		t.FailNow()
	}
	private , pub := keys.ExportJsonEncodedKeys()
	bprivate,err := json.Marshal(private)
	if err !=nil {
		t.Error("Serialize error",err)
		t.FailNow()
	}
	bpub,err := json.Marshal(pub)
	if err !=nil {
		t.Error("Serialize error",err)
		t.FailNow()
	}
	t.Log(string(bprivate))
	t.Log(string(bpub))
	msg := fimpgo.NewFloatMessage("evt.sensor.report", "temp_sensor", 35.5, nil, nil, nil)
	got, err := SignMessageES256(msg, nil, "alex@gmail.com", keys, nil)
	if err !=nil {
		t.Error("SerializedKey generation error",err)
		t.FailNow()
	}

	bmsg,err := got.SerializeToJson()
	if err !=nil {
		t.Error("Serialize error",err)
		t.FailNow()
	}

	t.Log(string(bmsg))

	keys2 := integration.NewEcdsaKey()
	err = keys2.ImportJsonPublicKey(pub)
	if err != nil {
		t.Error("Wrong key")
		t.FailNow()
	}
	result,err := GetVerifiedMessageES256(got,keys2)
	if err != nil || result == nil {
		t.Error("Signature is not valid")
	}else {
		t.Log(result)
	}
}