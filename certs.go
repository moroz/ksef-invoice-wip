package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"math/big"
	"time"
)

func GenerateSelfSignedCertificate(nip string) (derBytes []byte, key *ecdsa.PrivateKey, err error) {
	notBefore := time.Now().Add(-61 * time.Minute)
	notAfter := time.Now().AddDate(2, 0, 0)

	subject := pkix.Name{
		CommonName:   "JAN KOWALSKI",
		Country:      []string{"PL"},
		SerialNumber: "TINPL-" + nip,
		ExtraNames: []pkix.AttributeTypeAndValue{
			{
				Type:  asn1.ObjectIdentifier{2, 5, 4, 42},
				Value: "JAN",
			},
			{
				Type:  asn1.ObjectIdentifier{2, 5, 4, 4},
				Value: "KOWALSKI",
			},
		},
	}

	key, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return
	}

	serialBytes := make([]byte, 16)
	_, err = rand.Read(serialBytes)
	if err != nil {
		return
	}

	serialNumber := new(big.Int).SetBytes(serialBytes)

	template := x509.Certificate{
		Subject:               subject,
		SerialNumber:          serialNumber,
		NotAfter:              notAfter,
		NotBefore:             notBefore,
		KeyUsage:              x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
	}

	derBytes, err = x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	return
}
