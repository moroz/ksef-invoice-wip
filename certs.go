package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509/pkix"
)

func GenerateSelfSignedCertificate() error {
	// X509Certificate2 certificate = CertificateUtils.GetPersonalCertificate("A", "R", "TINPL", nip, "A R", EncryptionMethodEnum.ECDsa);

	parts := []string{
		"2.5.4.42=Jan",
		"2.5.4.4=Kowalski",
		"2.5.4.5=8976111986",
		"2.5.4.3=Jan Kowalski",
		"2.5.4.6=PL",
	}

	pkix.Name{
		Country: "PL",
	}

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}

}
