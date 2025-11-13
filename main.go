package main

import (
	"ksef-go/lib/certs"
	"ksef-go/lib/ksef"
	"log"
)

const SampleNip = "8976111986"

func main() {
	der, key, err := certs.GenerateSelfSignedCertificate(SampleNip)
	if err != nil {
		log.Fatal(err)
	}

	client := ksef.NewClient(SampleNip, der, key)

	if err := client.Authenticate(); err != nil {
		log.Fatal(err)
	}

	if client.Authenticated() {
		log.Print("Authenticated!")
	}
}
