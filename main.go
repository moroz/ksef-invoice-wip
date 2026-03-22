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
		log.Fatal("Load private key:", err)
	}

	der, err := certs.LoadKSeFCertificate(config.AuthenticationCertPath)
	if err != nil {
		log.Fatal("Load authentication certificate:", err)
	}

	client, err := ksef.NewClient(ksef.EnvironmentTest, config.NipNumber, der, key)
	if err != nil {
		log.Fatal("Initialize client:", err)
	}

	if err := client.Authenticate(); err != nil {
		log.Fatal("Authenticate:", err)
	}

	if client.Authenticated() {
		log.Print("Authenticated!")
	}

	session, err := client.OpenInteractiveSession()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%#v", session)
}
