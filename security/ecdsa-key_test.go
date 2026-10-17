package security

import "testing"

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
	if err := signer.ImportJsonPrivateKey(JsonEcKey{D: "0"}); err == nil {
		t.Fatal("expected error for zero scalar")
	}

	sig, err = SignStringES256("msg", signer)
	if err != nil {
		t.Fatal("Signing after failed import err:", err)
	}
	if !VerifyStringES256("msg", sig, verifier) {
		t.Fatal("Verification after failed import failed")
	}
}
