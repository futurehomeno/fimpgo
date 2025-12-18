package security

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"math/big"

	jwt "github.com/golang-jwt/jwt/v5"
)

type JsonEcKey struct {
	T string `json:"t"` //type - private/public
	X string `json:"x"`
	Y string `json:"y"`
	D string `json:"d"` // only for private key
}

type EcdsaKey struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

func (kp *EcdsaKey) SetPrivateKey(privateKey *ecdsa.PrivateKey) {
	kp.privateKey = privateKey
}

func (kp *EcdsaKey) PublicKey() *ecdsa.PublicKey {
	return kp.publicKey
}

func (kp *EcdsaKey) PrivateKey() *ecdsa.PrivateKey {
	return kp.privateKey
}

func NewEcdsaKey() *EcdsaKey {
	return &EcdsaKey{}
}

func (kp *EcdsaKey) Generate() error {
	var err error
	kp.privateKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}
	kp.publicKey = &kp.privateKey.PublicKey
	return nil
}

func (kp *EcdsaKey) ExportX509EncodedKeys() (string, string) {
	var pemEncodedStr, pemEncodedPubStr string
	if kp.privateKey != nil {
		x509Encoded, _ := x509.MarshalECPrivateKey(kp.privateKey)
		pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})
		pemEncodedStr = string(pemEncoded)
	}
	if kp.publicKey != nil {
		x509EncodedPub, _ := x509.MarshalPKIXPublicKey(kp.publicKey)
		pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})
		pemEncodedPubStr = string(pemEncodedPub)
	}
	return pemEncodedStr, pemEncodedPubStr
}

func (kp *EcdsaKey) ExportJsonEncodedKeys() (JsonEcKey, JsonEcKey) {
	privateKey := JsonEcKey{
		T: "private",
		X: kp.privateKey.X.Text(16),
		Y: kp.privateKey.Y.Text(16),
		D: kp.privateKey.D.Text(16),
	}
	pubKey := JsonEcKey{
		T: "public",
		X: kp.publicKey.X.Text(16),
		Y: kp.publicKey.Y.Text(16),
	}
	return privateKey, pubKey
}

func (kp *EcdsaKey) ImportX509PublicKey(pemEncodedPub string) error {
	blockPub, _ := pem.Decode([]byte(pemEncodedPub))
	if blockPub == nil {
		return errors.New("incorrect PEM format")
	}
	x509EncodedPub := blockPub.Bytes
	genericPublicKey, err := x509.ParsePKIXPublicKey(x509EncodedPub)
	if err != nil {
		return err
	}
	kp.publicKey = genericPublicKey.(*ecdsa.PublicKey)
	return nil
}

func (kp *EcdsaKey) ImportX509PrivateKey(pemEncoded string) error {
	var err error
	block, _ := pem.Decode([]byte(pemEncoded))
	x509Encoded := block.Bytes
	kp.privateKey, err = x509.ParseECPrivateKey(x509Encoded)
	if err != nil {
		return err
	}
	return nil
}

func (kp *EcdsaKey) ImportJsonPublicKey(jkey JsonEcKey) error {
	x, xok := big.NewInt(0).SetString(jkey.X, 16)
	y, yok := big.NewInt(0).SetString(jkey.Y, 16)
	if !xok || !yok {
		return errors.New("json key parse error")
	}
	kp.publicKey = &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}
	return nil
}

func (kp *EcdsaKey) ImportJsonPrivateKey(jkey JsonEcKey) error {
	x, xok := big.NewInt(0).SetString(jkey.X, 16)
	y, yok := big.NewInt(0).SetString(jkey.Y, 16)
	d, dok := big.NewInt(0).SetString(jkey.D, 16)
	if !xok || !yok || !dok {
		return errors.New("json key parse error")
	}
	kp.privateKey = &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: elliptic.P256(),
			X:     x,
			Y:     y,
		},
		D: d,
	}
	return nil
}

// SignStringES256 signs string and returns as result .
func SignStringES256(payload string, keys *EcdsaKey) (string, error) {
	signingMethodES256 := &jwt.SigningMethodECDSA{Name: "ES256", Hash: crypto.SHA256, KeySize: 32, CurveBits: 256}
	signature, err := signingMethodES256.Sign(payload, keys.PrivateKey())
	if err != nil {
		return "", err
	}
	return signature, nil
}

func VerifyStringES256(payload, sig string, key *EcdsaKey) bool {
	signingMethodES256 := &jwt.SigningMethodECDSA{Name: "ES256", Hash: crypto.SHA256, KeySize: 32, CurveBits: 256}
	err := signingMethodES256.Verify(payload, sig, key.PublicKey())
	if err == nil {
		return true
	} else {
		return false
	}
}
