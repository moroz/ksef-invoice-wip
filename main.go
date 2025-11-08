package main

import (
	"encoding/xml"
	"fmt"
	"log"
	xades "github.com/artemkunich/goxades"

)

const SampleNip = "8976111986"

func main() {
	client := &KSeFClient{}

	challenge, err := client.GetAuthChallenge()
	if err != nil {
		log.Fatal(err)
	}

	authTokenRequest := BuildAuthTokenRequestFromChallenge(challenge, SampleNip)
	output, err := xml.MarshalIndent(authTokenRequest, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(xml.Header, string(output))

	der, key, err := GenerateSelfSignedCertificate(SampleNip)
	if err != nil {
		log.Fatal(err)
	}

	keyStore := &xades.MemoryX509KeyStore{
		PrivateKey: key,
	}
}
