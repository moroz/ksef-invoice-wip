package ksef

const GetAuthEndpoint = "/auth/challenge"
const SubmitAuthXAdESSignatureEndpoint = "/auth/xades-signature"
const RedeemAuthTokenEndpoint = "/auth/token/redeem"
const CheckAuthenticationStatusEndpoint = "/auth/{ref}"

type Environment int

const (
	Environment_Test Environment = iota
	Environment_PreProd
	Envoronment_Prod
)

const BaseUrl_Test = "https://ksef-test.mf.gov.pl/api/v2"
const BaseUrl_PreProd = "https://ksef-demo.mf.gov.pl/api/v2"
