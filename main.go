package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"ksef-go/lib/certs"
	"ksef-go/lib/sign"

	"log"
)

const SampleNip = "8976111986"

func main() {
	der, key, err := certs.GenerateSelfSignedCertificate(SampleNip)
	if err != nil {
		log.Fatal(err)
	}

	client := &KSeFClient{}

	challenge, err := client.GetAuthChallenge()
	if err != nil {
		log.Fatal(err)
	}

	authTokenRequest := BuildAuthTokenRequestFromChallenge(challenge, SampleNip)
	xmlToSign, err := xml.Marshal(authTokenRequest)
	if err != nil {
		log.Fatal(err)
	}

	signed, err := sign.SignXML(string(xmlToSign), key, der)
	resp, err := client.SubmitAuthXadesSignature(string(signed))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(result))
}
