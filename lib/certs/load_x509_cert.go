package certs

import (
	"encoding/pem"
	"errors"
	"os"
)

func LoadKSeFCertificate(path string) (der []byte, err error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return
	}

	block, _ := pem.Decode(bytes)
	if block == nil {
		return nil, errors.New("failed to parse PEM")
	}

	return block.Bytes, nil
}
