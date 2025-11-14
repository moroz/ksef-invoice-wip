package ksef

import (
	"encoding/xml"
	"time"
)

type AuthChallengeResult struct {
	Challenge string
	Timestamp time.Time
}

type CheckAuthenticationStatusResponse struct {
	StartDate            time.Time
	AuthenticationMethod string
	Status               struct {
		Code        int
		Description string
	}
}

type RedeemAuthTokenResponse struct {
	AccessToken struct {
		Token      string
		ValidUntil time.Time
	}
	RefreshToken struct {
		Token      string
		ValidUntil time.Time
	}
}

// AuthTokenRequest represents the root element
type AuthTokenRequest struct {
	XMLName               xml.Name          `xml:"http://ksef.mf.gov.pl/auth/token/2.0 AuthTokenRequest"`
	XmlnsXsi              string            `xml:"xmlns:xsi,attr"`
	XmlnsXsd              string            `xml:"xmlns:xsd,attr"`
	Challenge             string            `xml:"Challenge"`
	ContextIdentifier     ContextIdentifier `xml:"ContextIdentifier"`
	SubjectIdentifierType string            `xml:"SubjectIdentifierType"`
}

type ContextIdentifier struct {
	Nip string `xml:"Nip"`
}

type SubmitAuthXAdESSignatureResponse struct {
	ReferenceNumber     string
	AuthenticationToken struct {
		Token      string
		ValidUntil time.Time
	}
}

type OpenInteractiveSessionRequest struct {
	FormCode   OpenInteractiveSessionRequestFormCode   `json:"formCode"`
	Encryption OpenInteractiveSessionRequestEncryption `json:"encryption"`
}

type OpenInteractiveSessionRequestFormCode struct {
	SystemCode    string `json:"systemCode"`
	SchemaVersion string `json:"schemaVersion"`
	Value         string `json:"value"`
}

type OpenInteractiveSessionRequestEncryption struct {
	EncryptedSymmetricKey []byte `json:"encryptedSymmetricKey"`
	InitializationVector  []byte `json:"initializationVector"`
}
