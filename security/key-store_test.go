package security

import (
	"path/filepath"
	"testing"
)

func TestKeyStore_GetEcdsaKey(t *testing.T) {
	user := "alex@mail.com"
	devId := "t-a-1"
	textToSign := "super important test message"
	path := filepath.Join(t.TempDir(), "key-store.json")
	store := NewKeyStore(path, false)
	ecKey := NewEcdsaKey()
	err := ecKey.Generate()
	if err != nil {
		t.Fatal("Generate err:", err)
	}

	_, public := ecKey.ExportX509EncodedKeys()

	sig, err := SignStringES256(textToSign, ecKey)
	if err != nil {
		t.Fatal("Signing err:", err)
	}
	// store.AddSerializedKey(user,devId,public,KeyTypePublic,AlgEcdsa256)
	if err := store.UpsertSerializedKey(user, devId, public, KeyTypePublic, AlgEcdsa256); err != nil {
		t.Fatal("Upserting err:", err)
	}

	store2 := NewKeyStore(path, false)
	if err := store2.LoadFromDisk(); err != nil {
		t.Fatal("LoadFromDisk err:", err)
	}

	keyFromStore, err := store2.GetEcdsaKey(user, devId, KeyTypePublic)
	if err != nil {
		t.Fatal("Get key err:", err)
	}

	isCorrect := VerifyStringES256(textToSign, sig, keyFromStore)
	if !isCorrect {
		t.Fatal("Verification failed")
	}
}
