package main

import (
	"ksef-go/lib/certs"
	"ksef-go/lib/ksef"
	"log"
)

func main() {
	der, key, err := certs.GenerateSelfSignedCertificate(SampleNip)
	if err != nil {
		log.Fatal(err)
	}

	client, err := ksef.NewClient(SampleNip, der, key)

	if err := client.Authenticate(); err != nil {
		log.Fatal(err)
	}

	if client.Authenticated() {
		log.Print("Authenticated!")
	}
}
