package main

import (
	_ "embed"
	"ksef-go/lib/certs"
	"ksef-go/lib/sign"
	"log"
	"os"

	"github.com/beevik/etree"
)

const SampleNip = "8976111986"

//go:embed formatted.xml
var example string

func main() {
	der, key, err := certs.GenerateSelfSignedCertificate(SampleNip)
	doc := etree.NewDocument()
	if err := doc.ReadFromString(example); err != nil {
		log.Fatal(err)
	}

	signature, err := sign.BuildXMLSignature(example, key, der)
	if err != nil {
		log.Fatal(err)
	}
	sigNode := etree.NewDocument()
	if err := sigNode.ReadFromBytes(signature); err != nil {
		log.Fatal(err)
	}

	doc.Root().AddChild(sigNode.Root())

	canonical, err := sign.CanonicalSerialize(doc.Root())
	if err != nil {
		log.Fatal(err)
	}

	os.Stdout.Write(canonical)
}
