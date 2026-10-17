package security

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type JsonEcKey struct {
	T string `json:"t"` // type - private/public
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

func (kp *EcdsaKey) ExportJsonEncodedKeys() (privKey JsonEcKey, pubKey JsonEcKey, err error) {
	// ----- PRIVATE KEY -----
	privBytes, err := kp.privateKey.Bytes()
	if err != nil {
		return privKey, pubKey, err
	}

	privateKey := JsonEcKey{
		T: "private",
		D: hex.EncodeToString(privBytes),
	}

	// ----- PUBLIC KEY -----
	pubBytes, err := kp.publicKey.Bytes()
	if err != nil {
		return privKey, pubKey, err
	}

	if len(pubBytes) == 0 || pubBytes[0] != 0x04 {
		return privKey, pubKey, errors.New("unexpected public key encoding")
	}

	coordLen := (len(pubBytes) - 1) / 2
	xBytes := pubBytes[1 : 1+coordLen]
	yBytes := pubBytes[1+coordLen:]

	privateKey.X = hex.EncodeToString(xBytes)
	privateKey.Y = hex.EncodeToString(yBytes)

	pubKey = JsonEcKey{
		T: "public",
		X: hex.EncodeToString(xBytes),
		Y: hex.EncodeToString(yBytes),
	}

	return privateKey, pubKey, nil
}

func (kp *EcdsaKey) ImportX509PublicKey(pemEncodedPub string) error {
	blockPub, _ := pem.Decode([]byte(pemEncodedPub))
	if blockPub == nil {
		return errors.New("incorrect PEM format")
	}
	genericPublicKey, err := x509.ParsePKIXPublicKey(blockPub.Bytes)
	if err != nil {
		return err
	}
	publicKey, ok := genericPublicKey.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("genericPublicKey cast type=%T fail", genericPublicKey)
	}
	kp.publicKey = publicKey
	return nil
}

func (kp *EcdsaKey) ImportX509PrivateKey(pemEncoded string) error {
	block, _ := pem.Decode([]byte(pemEncoded))
	if block == nil {
		return errors.New("incorrect PEM format")
	}
	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return err
	}
	kp.privateKey = privateKey
	return nil
}

func (kp *EcdsaKey) ImportJsonPublicKey(jkey JsonEcKey) error {
	x, xok := hexToFixed(jkey.X)
	y, yok := hexToFixed(jkey.Y)
	if !xok || !yok {
		return errors.New("json key parse error")
	}
	pub, err := ecdsa.ParseUncompressedPublicKey(elliptic.P256(), append(append([]byte{0x04}, x...), y...))
	if err != nil {
		return err
	}
	kp.publicKey = pub
	return nil
}

func (kp *EcdsaKey) ImportJsonPrivateKey(jkey JsonEcKey) error {
	d, ok := hexToFixed(jkey.D)
	if !ok {
		return errors.New("json key parse error")
	}
	priv, err := ecdsa.ParseRawPrivateKey(elliptic.P256(), d)
	if err != nil {
		return err
	}
	kp.privateKey = priv
	kp.publicKey = &priv.PublicKey
	return nil
}

// hexToFixed left-pads to the P-256 coordinate size, since keys exported via big.Int dropped leading zeros.
func hexToFixed(s string) ([]byte, bool) {
	n, ok := big.NewInt(0).SetString(s, 16)
	if !ok || strings.ContainsAny(s, "+-") || n.BitLen() > 256 {
		return nil, false
	}
	return n.FillBytes(make([]byte, 32)), true
}

// SignStringES256 signs string and returns as result .
func SignStringES256(payload string, keys *EcdsaKey) ([]byte, error) {
	signingMethodES256 := &jwt.SigningMethodECDSA{Name: "ES256", Hash: crypto.SHA256, KeySize: 32, CurveBits: 256}
	return signingMethodES256.Sign(payload, keys.PrivateKey())
}

func VerifyStringES256(payload string, sig []byte, key *EcdsaKey) bool {
	signingMethodES256 := &jwt.SigningMethodECDSA{Name: "ES256", Hash: crypto.SHA256, KeySize: 32, CurveBits: 256}
	err := signingMethodES256.Verify(payload, sig, key.PublicKey())
	return err == nil
}
