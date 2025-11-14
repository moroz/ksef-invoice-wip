package main

import (
	"encoding/pem"
	"fmt"
	"ksef-go/config"
	"log"
	"os"

	"github.com/youmark/pkcs8"
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

	key, err := pkcs8.ParsePKCS8PrivateKeyECDSA(block.Bytes, []byte(config.AuthenticationPrivKeyPassphrase))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", key)
}
