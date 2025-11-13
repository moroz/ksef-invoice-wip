package sign

import (
	"encoding/xml"
	"time"
)

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
