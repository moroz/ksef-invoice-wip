package ksef

import (
	_ "embed"
)

const GetAuthEndpoint = "/auth/challenge"
const SubmitAuthXAdESSignatureEndpoint = "/auth/xades-signature"
const RedeemAuthTokenEndpoint = "/auth/token/redeem"
const CheckAuthenticationStatusEndpoint = "/auth/{ref}"

type Environment int

const (
	EnvironmentTest Environment = iota
	EnvironmentDemo
	EnvironmentProd
)

const BaseurlTest = "https://api-test.ksef.mf.gov.pl/api/v2"
const BaseurlDemo = "https://api-demo.ksef.mf.gov.pl/api/v2"
const BaseurlProd = "https://api.ksef.mf.gov.pl/api/v2"

//go:embed certs/test.crt
var EncryptionCertTest []byte

//go:embed certs/prod.crt
var EncryptionCertProd []byte

//go:embed certs/demo.crt
var EncryptionCertDemo []byte

type EnvironmentConfig struct {
	EncryptionCert []byte
	BaseUrl        string
}

var EnvironmentConfigs = map[Environment]EnvironmentConfig{
	EnvironmentTest: {
		EncryptionCert: EncryptionCertTest,
		BaseUrl:        BaseurlTest,
	},
	EnvironmentDemo: {
		EncryptionCert: EncryptionCertDemo,
		BaseUrl:        BaseurlDemo,
	},
	EnvironmentProd: {
		EncryptionCert: EncryptionCertProd,
		BaseUrl:        BaseurlProd,
	},
}
