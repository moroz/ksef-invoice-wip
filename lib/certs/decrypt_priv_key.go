package certs

import (
	"crypto/ecdsa"
	"encoding/pem"
	"errors"
	"os"

	"github.com/youmark/pkcs8"
)

func DecryptKSeFPrivateKey(pemBytes, passphrase []byte) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("failed to decode private key from PEM")
	}

	return pkcs8.ParsePKCS8PrivateKeyECDSA(block.Bytes, passphrase)
}

func LoadKSeFPrivateKey(path, passphrase string) (*ecdsa.PrivateKey, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return DecryptKSeFPrivateKey(bytes, []byte(passphrase))
}
