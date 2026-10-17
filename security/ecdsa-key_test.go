package security

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"testing"
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
