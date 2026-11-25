package security

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeCertFiles(t *testing.T, dir string) {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(1), Subject: pkix.Name{CommonName: "test"},
		NotBefore: time.Now().Add(-time.Hour), NotAfter: time.Now().Add(time.Hour),
		IsCA: true, BasicConstraintsValid: true,
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	require.NoError(t, err)
	keyDer, err := x509.MarshalECPrivateKey(key)
	require.NoError(t, err)

	for name, block := range map[string]*pem.Block{
		"root-ca-1.pem": {Type: "CERTIFICATE", Bytes: der},
		"root-ca-2.pem": {Type: "CERTIFICATE", Bytes: der},
		"client.crt":    {Type: "CERTIFICATE", Bytes: der},
		"client.key":    {Type: "EC PRIVATE KEY", Bytes: keyDer},
	} {
		require.NoError(t, os.WriteFile(filepath.Join(dir, name), pem.EncodeToMemory(block), 0600))
	}
}

func TestTLSConfig(t *testing.T) {
	dir := t.TempDir()

	_, err := TLSConfig("", "", dir, false)
	assert.Error(t, err, "missing root CA")

	writeCertFiles(t, dir)

	cfg, err := TLSConfig("", "", dir, false)
	require.NoError(t, err)
	assert.NotNil(t, cfg.RootCAs)
	assert.Nil(t, cfg.NextProtos)
	assert.Empty(t, cfg.Certificates)
	assert.Equal(t, uint16(tls.VersionTLS12), cfg.MinVersion)

	cfg, err = TLSConfig("client.key", "client.crt", dir, true)
	require.NoError(t, err)
	assert.Equal(t, []string{"x-amzn-mqtt-ca"}, cfg.NextProtos)
	assert.Len(t, cfg.Certificates, 1)
	assert.NotNil(t, cfg.ClientCAs)
	assert.Equal(t, tls.RequireAndVerifyClientCert, cfg.ClientAuth)

	cfg, err = TLSConfig("", "client.crt", dir, false)
	require.NoError(t, err)
	assert.Empty(t, cfg.Certificates)

	for name, args := range map[string][2]string{
		"key without cert": {"client.key", ""},
		"missing cert":     {"", "nope.crt"},
		"missing key":      {"nope.key", "client.crt"},
	} {
		_, err = TLSConfig(args[0], args[1], dir, false)
		assert.Error(t, err, name)
	}
}

func TestCertPool_Missing(t *testing.T) {
	_, err := certPool(filepath.Join(t.TempDir(), "none"))
	assert.Error(t, err)
}
