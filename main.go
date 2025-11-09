package main

import (
	"crypto"
	"crypto/x509"
	"encoding/xml"
	"fmt"
	"io"

	xades "github.com/artemkunich/goxades"
	"github.com/beevik/etree"
	dsig "github.com/russellhaering/goxmldsig"

	"log"
)

const SampleNip = "8976111986"

func main() {
	der, key, err := GenerateSelfSignedRSACertificate(SampleNip)
	if err != nil {
		log.Fatal(err)
	}

	cert, err := x509.ParseCertificate(der)
	if err != nil {
		log.Fatal(err)
	}

	keyStore := &xades.MemoryX509KeyStore{
		PrivateKey: key,
		CertBinary: der,
		Cert:       cert,
	}

	client := &KSeFClient{}

	challenge, err := client.GetAuthChallenge()
	if err != nil {
		log.Fatal(err)
	}

	authTokenRequest := BuildAuthTokenRequestFromChallenge(challenge, SampleNip)
	output, err := xml.Marshal(authTokenRequest)
	if err != nil {
		log.Fatal(err)
	}

	doc := etree.NewDocument()
	err = doc.ReadFromString(xml.Header + string(output))
	root := removeComments(doc.Root())
	canonicalizer := dsig.MakeC14N10ExclusiveCanonicalizerWithPrefixList("")
	signContext := xades.SigningContext{
		DataContext: xades.SignedDataContext{
			Canonicalizer: canonicalizer,
			Hash:          crypto.SHA256,
			IsEnveloped:   true,
		},
		PropertiesContext: xades.SignedPropertiesContext{
			Canonicalizer: canonicalizer,
			Hash:          crypto.SHA256,
		},
		Canonicalizer: canonicalizer,
		Hash:          crypto.SHA256,
		KeyStore:      *keyStore,
	}

	signature, err := xades.CreateSignature(root, &signContext)
	if err != nil {
		log.Fatal(err)
	}

	root.AddChild(signature)

	b, err := canonicalSerialize(root)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.SubmitAuthXadesSignature(string(b))
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

func removeComments(elem *etree.Element) *etree.Element {
	copy := elem.Copy()
	for _, token := range copy.Child {
		_, ok := token.(*etree.Comment)
		if ok {
			copy.RemoveChild(token)
		}
	}
	for i, child := range copy.ChildElements() {
		copy.ChildElements()[i] = removeComments(child)
	}
	return copy
}

func canonicalSerialize(el *etree.Element) ([]byte, error) {
	doc := etree.NewDocument()
	doc.SetRoot(el.Copy())

	doc.WriteSettings = etree.WriteSettings{
		CanonicalAttrVal: true,
		CanonicalEndTags: false,
		CanonicalText:    true,
	}

	return doc.WriteToBytes()
}
