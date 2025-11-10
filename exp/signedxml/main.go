package main

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	_ "embed"
	"encoding/base64"
	"encoding/xml"
	"ksef-go/lib/certs"
	"log"
	"math/big"
	"os"

	"github.com/beevik/etree"
	dsig "github.com/russellhaering/goxmldsig"
)

const SampleNip = "8976111986"

//go:embed formatted.xml
var example string

// Signature represents the root XML-DSig Signature element
type Signature struct {
	XMLName xml.Name `xml:"http://www.w3.org/2000/09/xmldsig# Signature"`
	ID      string   `xml:"Id,attr,omitempty"`

	SignedInfo     SignedInfo     `xml:"SignedInfo"`
	SignatureValue SignatureValue `xml:"SignatureValue"`
}

// SignedInfo contains the canonicalization method, signature method, and references
type SignedInfo struct {
	XMLName                xml.Name               `xml:"SignedInfo"`
	CanonicalizationMethod CanonicalizationMethod `xml:"CanonicalizationMethod"`
	SignatureMethod        SignatureMethod        `xml:"SignatureMethod"`
	Reference              Reference              `xml:"Reference"`
}

// CanonicalizationMethod specifies the canonicalization algorithm
type CanonicalizationMethod struct {
	XMLName   xml.Name `xml:"CanonicalizationMethod"`
	Algorithm string   `xml:"Algorithm,attr"`
}

// SignatureMethod specifies the signature algorithm
type SignatureMethod struct {
	XMLName   xml.Name `xml:"SignatureMethod"`
	Algorithm string   `xml:"Algorithm,attr"`
}

// Reference contains information about the signed resource
type Reference struct {
	XMLName      xml.Name     `xml:"Reference"`
	URI          string       `xml:"URI,attr"`
	Transforms   Transforms   `xml:"Transforms"`
	DigestMethod DigestMethod `xml:"DigestMethod"`
	DigestValue  string       `xml:"DigestValue"`
}

// Transforms contains one or more Transform elements
type Transforms struct {
	XMLName   xml.Name    `xml:"Transforms"`
	Transform []Transform `xml:"Transform"`
}

// Transform specifies a transformation algorithm
type Transform struct {
	XMLName   xml.Name `xml:"Transform"`
	Algorithm string   `xml:"Algorithm,attr"`
}

// DigestMethod specifies the digest algorithm
type DigestMethod struct {
	XMLName   xml.Name `xml:"DigestMethod"`
	Algorithm string   `xml:"Algorithm,attr"`
}

// SignatureValue contains the actual signature
type SignatureValue struct {
	XMLName xml.Name `xml:"SignatureValue"`
	Value   string   `xml:",chardata"`
}

func SerializeSignature(key *ecdsa.PrivateKey, r, s *big.Int) []byte {
	byteCount := (key.Curve.Params().BitSize + 7) / 8
	rBytes := r.FillBytes(make([]byte, byteCount))
	sBytes := r.FillBytes(make([]byte, byteCount))
	return append(rBytes, sBytes...)
}

func main() {
	_, key, err := certs.GenerateSelfSignedCertificate(SampleNip)

	excC14N := dsig.MakeC14N10ExclusiveCanonicalizerWithPrefixList("")

	doc := etree.NewDocument()
	if err := doc.ReadFromString(example); err != nil {
		log.Fatal(err)
	}

	canonical, err := excC14N.Canonicalize(doc.Root())
	if err != nil {
		log.Fatal(err)
	}

	digestBytes := sha256.Sum256(canonical)
	digestValue := base64.StdEncoding.EncodeToString(digestBytes[:])

	sig := &Signature{
		ID: "Signature",
		SignedInfo: SignedInfo{
			CanonicalizationMethod: CanonicalizationMethod{
				Algorithm: "http://www.w3.org/TR/2001/REC-xml-c14n-20010315",
			},
			SignatureMethod: SignatureMethod{
				Algorithm: "http://www.w3.org/2001/04/xmldsig-more#ecdsa-sha256",
			},
			Reference: Reference{
				Transforms: Transforms{
					Transform: []Transform{
						{
							Algorithm: "http://www.w3.org/2000/09/xmldsig#enveloped-signature",
						},
						{
							Algorithm: "http://www.w3.org/2001/10/xml-exc-c14n#",
						},
					},
				},
				DigestMethod: DigestMethod{
					Algorithm: "http://www.w3.org/2001/04/xmlenc#sha256",
				},
				DigestValue: digestValue,
			},
		},
		SignatureValue: SignatureValue{},
	}

	signedInfoXml, err := xml.Marshal(sig.SignedInfo)
	if err != nil {
		log.Fatal(err)
	}

	parsedSignedInfo := etree.NewDocument()
	if err := parsedSignedInfo.ReadFromBytes(signedInfoXml); err != nil {
		log.Fatal(err)
	}

	canonical, err = dsig.MakeC14N10RecCanonicalizer().Canonicalize(parsedSignedInfo.Root())
	if err != nil {
		log.Fatal(err)
	}

	signedInfoHash := sha256.Sum256(canonical)
	r, s, err := ecdsa.Sign(rand.Reader, key, signedInfoHash[:])
	if err != nil {
		log.Fatal(err)
	}
	signatureBytes := SerializeSignature(key, r, s)

	sig.SignatureValue.Value = base64.StdEncoding.EncodeToString(signatureBytes)

	serialized, err := xml.Marshal(sig)
	sigNode := etree.NewDocument()
	if err := sigNode.ReadFromBytes(serialized); err != nil {
		log.Fatal(err)
	}

	doc.Root().AddChild(sigNode.Root())

	canonical, err = canonicalSerialize(doc.Root())
	if err != nil {
		log.Fatal(err)
	}

	os.Stdout.Write(canonical)
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
