package main

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"ksef-go/config"
	"log"
	"os"
)

func main() {
	pemBytes, err := os.ReadFile(config.AuthenticationPrivKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	block, _ := pem.Decode(pemBytes)
	if block == nil {
		log.Fatal("Failed to decode priv key PEM")
	}

	decrypted, err := x509.DecryptPEMBlock(block, []byte(config.AuthenticationPrivKeyPassphrase))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(decrypted)
}
