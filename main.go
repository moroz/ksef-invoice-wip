package main

import (
	"ksef-go/config"
	"ksef-go/lib/certs"
	"ksef-go/lib/ksef"
	"log"
)

func main() {
	key, err := certs.LoadKSeFPrivateKey(config.AuthenticationPrivKeyPath, config.AuthenticationPrivKeyPassphrase)
	if err != nil {
		log.Fatal(err)
	}

	der, err := certs.LoadKSeFCertificate(config.AuthenticationCertPath)
	if err != nil {
		log.Fatal(err)
	}

	client, err := ksef.NewClient(ksef.EnvironmentProd, config.NipNumber, der, key)
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Authenticate(); err != nil {
		log.Fatal(err)
	}

	if client.Authenticated() {
		log.Print("Authenticated!")
	}
}
