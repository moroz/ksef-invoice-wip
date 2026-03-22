package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"ksef-go/lib/ksef"
	"log"
)

func main() {
	block, _ := pem.Decode(ksef.EncryptionCertProd)
	cert, _ := x509.ParseCertificate(block.Bytes)

	if key, ok := cert.PublicKey.(*rsa.PublicKey); !ok {
		log.Fatalf("Not ok")
	} else {
		fmt.Printf("%#v\n", key)
	}

}
