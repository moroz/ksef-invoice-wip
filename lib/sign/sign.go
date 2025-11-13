package sign

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/xml"
	"math/big"
	"time"

	"github.com/beevik/etree"
	dsig "github.com/russellhaering/goxmldsig"
)

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
					Algorithm: AlgSHA256,
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
				Algorithm: ExclusiveC14N,
			},
			SignatureMethod: SignatureMethod{
				Algorithm: AlgECDSA_SHA256,
			},
			Reference: []Reference{
				{
					Transforms: Transforms{
						Transform: []Transform{
							{
								Algorithm: EnvelopedSignature,
							},
							{
								Algorithm: ExclusiveC14N,
							},
						},
					},
					DigestMethod: DigestMethod{
						Algorithm: AlgSHA256,
					},
				},
				{
					URI:  "#SignedProperties",
					Type: ETSISignedProperties,
					Transforms: Transforms{
						Transform: []Transform{
							{
								Algorithm: ExclusiveC14N,
							},
						},
					},
					DigestMethod: DigestMethod{
						Algorithm: AlgSHA256,
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
				XMLNS:  ETSIQualifyingProperties,
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

func SignXML(xmlString string, key *ecdsa.PrivateKey, certDer []byte) (result []byte, err error) {
	doc := etree.NewDocument()
	if err = doc.ReadFromString(xmlString); err != nil {
		return
	}

	signature, err := BuildXMLSignature(xmlString, key, certDer)
	if err != nil {
		return
	}

	sigNode := etree.NewDocument()
	if err = sigNode.ReadFromBytes(signature); err != nil {
		return
	}

	doc.Root().AddChild(sigNode.Root())

	return CanonicalSerialize(doc.Root())
}
