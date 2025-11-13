package main

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	_ "embed"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"ksef-go/lib/certs"
	"log"
	"math/big"
	"os"
	"time"

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
	KeyInfo        KeyInfo        `xml:"KeyInfo"`
	Object         XAdESObject    `xml:"Object"`
}

type KeyInfo struct {
	X509Data X509Data `xml:"X509Data"`
}

type X509Data struct {
	X509Certificates []X509Certificate `xml:"X509Certificate"`
}

type X509Certificate struct {
	XMLName xml.Name `xml:"X509Certificate"`
	Value   string   `xml:",chardata"`
}

// SignedInfo contains the canonicalization method, signature method, and references
type SignedInfo struct {
	XMLName                xml.Name               `xml:"http://www.w3.org/2000/09/xmldsig# SignedInfo"`
	CanonicalizationMethod CanonicalizationMethod `xml:"CanonicalizationMethod"`
	SignatureMethod        SignatureMethod        `xml:"SignatureMethod"`
	Reference              []Reference            `xml:"Reference"`
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
	Type         string       `xml:"Type,attr,omitempty"`
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

type XAdESObject struct {
	XMLName              xml.Name `xml:"Object"`
	QualifyingProperties XAdESQualifyingProperties
}

type XAdESQualifyingProperties struct {
	XMLName          xml.Name `xml:"xades:QualifyingProperties"`
	XMLNS            string   `xml:"xmlns:xades,attr"`
	Target           string   `xml:"Target,attr"`
	SignedProperties XAdESSignedProperties
}

type XAdESSignedProperties struct {
	XMLName                   xml.Name `xml:"xades:SignedProperties"`
	ID                        string   `xml:"Id,attr"`
	SignedSignatureProperties XAdESSignedSignatureProperties
}

type XAdESSignedSignatureProperties struct {
	XMLName            xml.Name  `xml:"xades:SignedSignatureProperties"`
	SigningTime        time.Time `xml:"xades:SigningTime"`
	SigningCertificate XAdESSigningCertificate
}

type XAdESSigningCertificate struct {
	XMLName xml.Name `xml:"xades:SigningCertificate"`
	Cert    XAdESCert
}

type XAdESCert struct {
	XMLName      xml.Name `xml:"xades:Cert"`
	CertDigest   XAdESCertDigest
	IssuerSerial XAdESIssuerSerial
}

type XAdESCertDigest struct {
	XMLName      xml.Name     `xml:"xades:CertDigest"`
	DigestMethod DigestMethod `xml:"DigestMethod"`
	DigestValue  string       `xml:"DigestValue"`
}

type XAdESIssuerSerial struct {
	XMLName          xml.Name `xml:"xades:IssuerSerial"`
	X509IssuerName   string   `xml:"X509IssuerName"`
	X509SerialNumber string   `xml:"X509SerialNumber"`
}

func SerializeSignatureXMLDSig(r, s *big.Int) string {
	// For P-256 (secp256r1), R and S should each be 32 bytes
	rBytes := r.Bytes()
	sBytes := s.Bytes()

	// Pad to 32 bytes if needed
	signature := make([]byte, 64)
	copy(signature[32-len(rBytes):32], rBytes)
	copy(signature[64-len(sBytes):64], sBytes)

	return base64.StdEncoding.EncodeToString(signature)
}

func buildSigningCertificateInfo(certDer []byte) (*XAdESSigningCertificate, error) {
	cert, err := x509.ParseCertificate(certDer)
	if err != nil {
		return nil, err
	}

	certHash := sha256.Sum256(certDer)

	return &XAdESSigningCertificate{
		Cert: XAdESCert{
			CertDigest: XAdESCertDigest{
				DigestMethod: DigestMethod{
					Algorithm: "http://www.w3.org/2001/04/xmlenc#sha256",
				},
				DigestValue: base64.StdEncoding.EncodeToString(certHash[:]),
			},
			IssuerSerial: XAdESIssuerSerial{
				X509IssuerName:   cert.Issuer.String(),
				X509SerialNumber: cert.SerialNumber.String(),
			},
		},
	}, nil
}

func BuildSignatureTemplate(certDer []byte) (*Signature, error) {
	certInfo, err := buildSigningCertificateInfo(certDer)
	if err != nil {
		return nil, err
	}

	return &Signature{
		ID: "Signature",
		SignedInfo: SignedInfo{
			CanonicalizationMethod: CanonicalizationMethod{
				Algorithm: "http://www.w3.org/2001/10/xml-exc-c14n#",
			},
			SignatureMethod: SignatureMethod{
				Algorithm: "http://www.w3.org/2001/04/xmldsig-more#ecdsa-sha256",
			},
			Reference: []Reference{
				{
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
				},
				{
					URI:  "#SignedProperties",
					Type: "http://uri.etsi.org/01903#SignedProperties",
					Transforms: Transforms{
						Transform: []Transform{
							{
								Algorithm: "http://www.w3.org/2001/10/xml-exc-c14n#",
							},
						},
					},
					DigestMethod: DigestMethod{
						Algorithm: "http://www.w3.org/2001/04/xmlenc#sha256",
					},
				},
			},
		},
		SignatureValue: SignatureValue{},
		KeyInfo: KeyInfo{X509Data: X509Data{
			X509Certificates: []X509Certificate{
				{
					Value: base64.StdEncoding.EncodeToString(certDer),
				},
			},
		}},
		Object: XAdESObject{
			QualifyingProperties: XAdESQualifyingProperties{
				XMLNS:  "http://uri.etsi.org/01903/v1.3.2#",
				Target: "#Signature",
				SignedProperties: XAdESSignedProperties{
					ID: "SignedProperties",
					SignedSignatureProperties: XAdESSignedSignatureProperties{
						SigningTime:        time.Now(),
						SigningCertificate: *certInfo,
					},
				},
			},
		},
	}, nil
}

func MarshalToEtreeDocument(object any) (*etree.Document, error) {
	bytes, err := xml.Marshal(object)
	if err != nil {
		return nil, err
	}

	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(bytes); err != nil {
		return nil, err
	}
	return doc, nil
}

func calculateDigest(template *Signature, path string, c14n dsig.Canonicalizer) ([]byte, error) {
	document, err := MarshalToEtreeDocument(template)
	if err != nil {
		return nil, err
	}

	element := document.Root().FindElement(path)
	if element == nil {
		return nil, fmt.Errorf("calculateDigest: element with path %s not found in root", path)
	}

	canonical, err := c14n.Canonicalize(element)
	if err != nil {
		return nil, err
	}

	//fmt.Println(string(canonical))

	digest := sha256.Sum256(canonical)

	return digest[:], nil
}

func calculateSignature(template *Signature, path string, key *ecdsa.PrivateKey, c14n dsig.Canonicalizer) (*big.Int, *big.Int, error) {
	digest, err := calculateDigest(template, path, c14n)
	if err != nil {
		return nil, nil, err
	}
	return ecdsa.Sign(rand.Reader, key, digest)
}

func buildQualifiedSignedProperties(signedProperties *etree.Element) *etree.Element {
	qualifiedSignedProperties := signedProperties.Copy()
	qualifiedSignedProperties.Attr = append(
		signedProperties.Attr,
		etree.Attr{Space: "xmlns", Key: "", Value: dsig.Namespace},
		etree.Attr{Space: "xmlns", Key: "xades", Value: "http://uri.etsi.org/01903/v1.3.2#"},
	)
	return qualifiedSignedProperties
}

func calculateSignedPropertiesDigest(signature *Signature, c14n dsig.Canonicalizer) ([]byte, error) {
	doc, err := MarshalToEtreeDocument(signature)
	if err != nil {
		return nil, err
	}

	element := doc.Root().FindElement(".//[@Id='SignedProperties']")
	if element == nil {
		return nil, fmt.Errorf("SignedProperties element not found")
	}

	qualified := buildQualifiedSignedProperties(element)

	canonical, err := c14n.Canonicalize(qualified)
	if err != nil {
		return nil, err
	}

	digest := sha256.Sum256(canonical)
	return digest[:], nil
}

func BuildXMLSignature(xmlString string, key *ecdsa.PrivateKey, certDer []byte) ([]byte, error) {
	c14n := dsig.MakeC14N10ExclusiveCanonicalizerWithPrefixList("")

	doc := etree.NewDocument()
	if err := doc.ReadFromString(xmlString); err != nil {
		return nil, err
	}

	canonicalRoot, err := c14n.Canonicalize(doc.Root())
	if err != nil {
		return nil, err
	}

	digestBytes := sha256.Sum256(canonicalRoot)

	signature, err := BuildSignatureTemplate(certDer)
	if err != nil {
		return nil, err
	}
	signature.SignedInfo.Reference[0].DigestValue = base64.StdEncoding.EncodeToString(digestBytes[:])

	signedInfoDigest, err := calculateSignedPropertiesDigest(signature, c14n)
	if err != nil {
		return nil, err
	}
	signature.SignedInfo.Reference[1].DigestValue = base64.StdEncoding.EncodeToString(signedInfoDigest)

	r, s, err := calculateSignature(signature, "./SignedInfo", key, c14n)
	if err != nil {
		return nil, err
	}
	signature.SignatureValue.Value = SerializeSignatureXMLDSig(r, s)
	return xml.Marshal(signature)
}

func main() {
	der, key, err := certs.GenerateSelfSignedCertificate(SampleNip)
	doc := etree.NewDocument()
	if err := doc.ReadFromString(example); err != nil {
		log.Fatal(err)
	}

	signature, err := BuildXMLSignature(example, key, der)
	if err != nil {
		log.Fatal(err)
	}
	sigNode := etree.NewDocument()
	if err := sigNode.ReadFromBytes(signature); err != nil {
		log.Fatal(err)
	}

	doc.Root().AddChild(sigNode.Root())

	canonical, err := canonicalSerialize(doc.Root())
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
