package security

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

func caCertPool(certDir string) (*x509.CertPool, error) {
	certs := x509.NewCertPool()
	cafile := filepath.Join(certDir, "root-ca-1.pem")
	pemData, err := os.ReadFile(cafile)
	if err != nil {
		return nil, err
	}
	certs.AppendCertsFromPEM(pemData)

	// Optional additional CAs
	for _, ca := range []string{"root-ca-2.pem", "root-ca-3.pem"} {
		cafile = filepath.Join(certDir, ca)
		if pemData, err = os.ReadFile(cafile); err == nil {
			certs.AppendCertsFromPEM(pemData)
		}
	}

	log.Infof("[fimpgo] CA certificates are loaded")
	return certs, nil
}

func certPool(certFile string) (*x509.CertPool, error) {
	certs := x509.NewCertPool()
	pemData, err := os.ReadFile(certFile)
	if err != nil {
		return nil, err
	}
	certs.AppendCertsFromPEM(pemData)
	log.Infof("[fimpgo] Certificate %v is loaded", certFile)
	return certs, nil
}

// ConfigureTLS method should be used to configure mutual TLS like AwS IoT core does. Also it configures TLS protocol switch .
// Cert dir should contains all CA root certificates .
// IsAws flag controls AWS specific TLS protocol switch.
func TLSConfig(privateKeyFileName, certFileName, certDir string, isAWS bool) (*tls.Config, error) {
	hasCert := certFileName != ""
	hasKey := privateKeyFileName != ""

	privateKeyFileName = filepath.Join(certDir, privateKeyFileName)
	certFileName = filepath.Join(certDir, certFileName)
	config := &tls.Config{InsecureSkipVerify: false}

	if isAWS {
		config.NextProtos = []string{"x-amzn-mqtt-ca"}
	}

	var err error
	config.RootCAs, err = caCertPool(certDir)
	if err != nil {
		return nil, err
	}

	if hasCert {
		config.ClientCAs, err = certPool(certFileName)
		if err != nil {
			return nil, err
		}
		config.ClientAuth = tls.RequireAndVerifyClientCert
	}

	if hasKey {
		if !hasCert {
			return nil, fmt.Errorf("key specified but cert is not")
		}

		cert, err := tls.LoadX509KeyPair(certFileName, privateKeyFileName)
		if err != nil {
			return nil, err
		}

		config.Certificates = []tls.Certificate{cert}
	}

	return config, nil
}
