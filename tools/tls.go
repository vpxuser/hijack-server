package tools

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/vpxuser/proxy"
	"os"
)

func decodePEM(path string) (*pem.Block, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		proxy.Fatal(err)
	}

	block, _ := pem.Decode(file)
	if block == nil || len(block.Bytes) <= 0 {
		return nil, errors.New("failed to decode PEM block")
	}

	return block, nil
}

func LoadCert(path string, cert **x509.Certificate) {
	block, err := decodePEM(path)
	if err != nil {
		proxy.Fatal(err)
	}

	*cert, err = x509.ParseCertificate(block.Bytes)
	if err != nil {
		proxy.Fatal(err)
	}
}

func LoadKey(path string, key *crypto.PrivateKey) {
	block, err := decodePEM(path)
	if err != nil {
		proxy.Fatal(err)
	}

	switch block.Type {
	case "EC PRIVATE KEY":
		*key, err = x509.ParseECPrivateKey(block.Bytes)
	case "RSA PRIVATE KEY":
		*key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	case "PRIVATE KEY":
		*key, err = x509.ParsePKCS8PrivateKey(block.Bytes)
	}

	if err != nil {
		proxy.Fatal(err)
	}
}
