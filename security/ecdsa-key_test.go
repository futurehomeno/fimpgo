package security

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEcdsaKey_JsonRoundTrip(t *testing.T) {
	key := NewEcdsaKey()
	if err := key.Generate(); err != nil {
		t.Fatal("Generate err:", err)
	}
	priv, pub, err := key.ExportJsonEncodedKeys()
	if err != nil {
		t.Fatal("Export err:", err)
	}

	signer := NewEcdsaKey()
	if err := signer.ImportJsonPrivateKey(priv); err != nil {
		t.Fatal("Import private err:", err)
	}
	verifier := NewEcdsaKey()
	if err := verifier.ImportJsonPublicKey(pub); err != nil {
		t.Fatal("Import public err:", err)
	}

	sig, err := SignStringES256("msg", signer)
	if err != nil {
		t.Fatal("Signing err:", err)
	}
	if !VerifyStringES256("msg", sig, verifier) {
		t.Fatal("Verification failed")
	}
	if err := verifier.ImportJsonPublicKey(JsonEcKey{X: "01", Y: "02"}); err == nil {
		t.Fatal("expected error for point off the curve")
	}
	if err := verifier.ImportJsonPublicKey(JsonEcKey{X: "-" + pub.X, Y: pub.Y}); err == nil {
		t.Fatal("expected error for signed coordinate")
	}
	if err := signer.ImportJsonPrivateKey(JsonEcKey{D: "-1"}); err == nil {
		t.Fatal("expected error for signed scalar")
	}
	if err := signer.ImportJsonPrivateKey(JsonEcKey{D: "+1"}); err == nil {
		t.Fatal("expected error for plus-signed scalar")
	}
}

func TestEcdsaKey_FailedX509ImportKeepsKey(t *testing.T) {
	key := NewEcdsaKey()
	if err := key.Generate(); err != nil {
		t.Fatal("Generate err:", err)
	}
	edPub, _, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	edDer, err := x509.MarshalPKIXPublicKey(edPub)
	if err != nil {
		t.Fatal(err)
	}

	if err := key.ImportX509PublicKey(string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: edDer}))); err == nil {
		t.Error("expected error for non-ECDSA public key")
	}
	if err := key.ImportX509PrivateKey("not pem"); err == nil {
		t.Error("expected error for non-PEM private key")
	}
	if err := key.ImportX509PrivateKey(string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: []byte{1, 2, 3}}))); err == nil {
		t.Error("expected error for unparsable private key")
	}

	sig, err := SignStringES256("msg", key)
	if err != nil {
		t.Fatal("Signing err:", err)
	}
	if !VerifyStringES256("msg", sig, key) {
		t.Fatal("Verification failed")
	}
}

func TestEcdsaKey_X509RoundTripAndAccessors(t *testing.T) {
	key := NewEcdsaKey()
	priv, pub := key.ExportX509EncodedKeys()
	assert.Empty(t, priv)
	assert.Empty(t, pub)
	require.NoError(t, key.Generate())
	priv, pub = key.ExportX509EncodedKeys()

	other := NewEcdsaKey()
	require.NoError(t, other.ImportX509PrivateKey(priv))
	require.NoError(t, other.ImportX509PublicKey(pub))
	assert.True(t, key.PrivateKey().Equal(other.PrivateKey()))
	assert.True(t, key.PublicKey().Equal(other.PublicKey()))

	third := NewEcdsaKey()
	third.SetPrivateKey(key.PrivateKey())
	assert.Equal(t, key.PrivateKey(), third.PrivateKey())

	assert.Error(t, other.ImportX509PublicKey("junk"))
	assert.Error(t, other.ImportX509PublicKey(string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: []byte{1}}))))
}

func TestEcdsaKey_JsonImportErrors(t *testing.T) {
	k := NewEcdsaKey()
	assert.Error(t, k.ImportJsonPublicKey(JsonEcKey{X: "zz", Y: "01"}))
	assert.Error(t, k.ImportJsonPublicKey(JsonEcKey{X: "01", Y: "zz"}))
	assert.Error(t, k.ImportJsonPrivateKey(JsonEcKey{D: "zz"}))
	assert.Error(t, k.ImportJsonPrivateKey(JsonEcKey{D: "0"}))
	assert.Error(t, k.ImportJsonPrivateKey(JsonEcKey{D: "1" + strings.Repeat("0", 64)}))
}

func TestVerifyStringES256_Negative(t *testing.T) {
	key := NewEcdsaKey()
	require.NoError(t, key.Generate())
	other := NewEcdsaKey()
	require.NoError(t, other.Generate())
	sig, err := SignStringES256("msg", key)
	require.NoError(t, err)
	assert.False(t, VerifyStringES256("tampered", sig, key))
	assert.False(t, VerifyStringES256("msg", sig, other))
	assert.False(t, VerifyStringES256("msg", sig[:10], key))
}
