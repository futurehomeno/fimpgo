package transport

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
)


type EcKeyPair struct {
	privateKey *ecdsa.PrivateKey
	publicKey *ecdsa.PublicKey
}

func (kp *EcKeyPair) SetPrivateKey(privateKey *ecdsa.PrivateKey) {
	kp.privateKey = privateKey
}

func (kp *EcKeyPair) PublicKey() *ecdsa.PublicKey {
	return kp.publicKey
}

func (kp *EcKeyPair) PrivateKey() *ecdsa.PrivateKey {
	return kp.privateKey
}

func NewEcKeyPair() *EcKeyPair {
	return &EcKeyPair{}
}

func (kp *EcKeyPair) Generate() error {
	var err error
	kp.privateKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}
	kp.publicKey = &kp.privateKey.PublicKey
	return nil
}

func (kp *EcKeyPair) GetEncodedKeys() (string,string) {
	x509Encoded, _ := x509.MarshalECPrivateKey(kp.privateKey)
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})

	x509EncodedPub, _ := x509.MarshalPKIXPublicKey(kp.publicKey)
	pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})
	return string(pemEncoded), string(pemEncodedPub)
}

func (kp *EcKeyPair) ImportPublicKey(pemEncodedPub string)  error {
	blockPub, _ := pem.Decode([]byte(pemEncodedPub))
	x509EncodedPub := blockPub.Bytes
	genericPublicKey, err := x509.ParsePKIXPublicKey(x509EncodedPub)
	if err != nil {
		return err
	}
	kp.publicKey = genericPublicKey.(*ecdsa.PublicKey)
	return  nil
}

func (kp *EcKeyPair) ImportPrivateKey(pemEncoded string) error {
	var err error
	block, _ := pem.Decode([]byte(pemEncoded))
	x509Encoded := block.Bytes
	kp.privateKey, err = x509.ParseECPrivateKey(x509Encoded)
	if err != nil {
		return err
	}
	return nil
}

