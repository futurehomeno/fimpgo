package security

import (
	"testing"
)

func TestKeyStore_GetEcdsaKey(t *testing.T) {
	user := "alex@mail.com"
	devId := "t-a-1"
	textToSign := "super important test message"
	store := NewKeyStore("test-key-store.json",false)
	ecKey := NewEcdsaKey()
	ecKey.Generate()

	_,public := ecKey.ExportX509EncodedKeys()

	sig , err := SignStringES256(textToSign,ecKey)
	if err != nil {
		t.Fatal("signing error . Err:",err)
	}
	store.AddSerializedKey(user,devId,public,KeyTypePublic,AlgEcdsa256)

	store2 := NewKeyStore("test-key-store.json",false)
	store2.LoadFromDisk()
	keyFromStore,err := store2.GetEcdsaKey(user,devId,KeyTypePublic)
	if err != nil {
		t.Fatal("can't get the key . Err:",err)
	}

	isCorrect := VerifyStringES256(textToSign,sig,keyFromStore)
	if !isCorrect {
		t.Fatal("Verification failed")
	}
}