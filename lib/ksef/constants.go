package ksef

const GetAuthEndpoint = "/auth/challenge"
const SubmitAuthXAdESSignatureEndpoint = "/auth/xades-signature"
const RedeemAuthTokenEndpoint = "/auth/token/redeem"
const CheckAuthenticationStatusEndpoint = "/auth/{ref}"

type Environment int

const (
	EnvironmentTest Environment = iota
	EnvironmentPreprod
	EnvironmentProd
)

const BaseurlTest = "https://ksef-test.mf.gov.pl/api/v2"
const BaseurlPreprod = "https://ksef-demo.mf.gov.pl/api/v2"
const BaseurlProd = "https://ksef.mf.gov.pl/api/v2"
