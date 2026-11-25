package security

import (
	"os"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestKeyStore_Records(t *testing.T) {
	path := filepath.Join(t.TempDir(), "ks.json")
	s := NewKeyStore(path, true)
	require.NoError(t, s.AddSerializedKey("u1", "d1", "k1", KeyTypePublic, AlgEcdsa256))
	require.NoError(t, s.AddSerializedKey("u1", "d2", "k2", KeyTypeSymmetric, "HS"))
	assert.Equal(t, "/var/lib/futurehome/hub/pub_key_store.json", NewKeyStore("", false).keyStoreFilePath)

	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0600), info.Mode().Perm())

	assert.True(t, s.CheckIfUserHasKey("u1", "d1"))
	assert.False(t, s.CheckIfUserHasKey("u1", "d3"))
	assert.Len(t, s.GetAllUserKeys("u1"), 2)
	assert.Empty(t, s.GetAllUserKeys("u2"))

	assert.Equal(t, "k2", s.GetKey("u1", "d2", KeyTypeSymmetric).SerializedKey)
	assert.Nil(t, s.GetKey("u1", "d2", KeyTypePublic))

	updated, err := s.UpdateSerializedKey("u1", "d1", "k1b", KeyTypePublic, AlgEcdsa256)
	require.NoError(t, err)
	assert.True(t, updated)
	updated, err = s.UpdateSerializedKey("u9", "d1", "x", KeyTypePublic, AlgEcdsa256)
	require.NoError(t, err)
	assert.False(t, updated)

	require.NoError(t, s.UpsertSerializedKey("u1", "d1", "k1c", KeyTypePublic, AlgEcdsa256))
	require.NoError(t, s.UpsertSerializedKey("u3", "d1", "k3", KeyTypePublic, AlgEcdsa256))
	assert.Equal(t, "k1c", s.GetKey("u1", "d1", KeyTypePublic).SerializedKey)
	assert.Len(t, s.keyStore, 3)

	s2 := NewKeyStore(path, true)
	require.NoError(t, s2.LoadFromDisk())
	assert.Len(t, s2.keyStore, 3)

	assert.Error(t, NewKeyStore(filepath.Join(t.TempDir(), "none"), false).LoadFromDisk())
	require.NoError(t, os.WriteFile(path, []byte("{"), 0600))
	assert.Error(t, s2.LoadFromDisk())
	assert.Error(t, NewKeyStore(filepath.Join(t.TempDir(), "no", "ks"), false).SaveToDisk())
}

func TestKeyStore_GetEcdsaKeyErrors(t *testing.T) {
	s := NewKeyStore(filepath.Join(t.TempDir(), "ks.json"), false)
	_, err := s.GetEcdsaKey("u", "d", KeyTypePublic)
	assert.EqualError(t, err, "key not found")

	require.NoError(t, s.AddSerializedKey("u", "d", "", KeyTypePublic, AlgEcdsa256))
	_, err = s.GetEcdsaKey("u", "d", KeyTypePublic)
	assert.EqualError(t, err, "empty key string")

	require.NoError(t, s.AddSerializedKey("u", "d", "garbage", KeyTypePrivate, AlgEcdsa256))
	_, err = s.GetEcdsaKey("u", "d", KeyTypePrivate)
	assert.Error(t, err)

	require.NoError(t, s.AddSerializedKey("u", "d", "x", KeyTypeSymmetric, AlgEcdsa256))
	_, err = s.GetEcdsaKey("u", "d", KeyTypeSymmetric)
	assert.Error(t, err)

	key := NewEcdsaKey()
	require.NoError(t, key.Generate())
	priv, _ := key.ExportX509EncodedKeys()
	require.NoError(t, s.AddSerializedKey("u2", "d", priv, KeyTypePrivate, AlgEcdsa256))
	got, err := s.GetEcdsaKey("u2", "d", KeyTypePrivate)
	require.NoError(t, err)
	again, err := s.GetEcdsaKey("u2", "d", KeyTypePrivate)
	require.NoError(t, err)
	assert.Same(t, got, again)
	assert.Equal(t, key.PrivateKey(), got.PrivateKey())
}
