package main

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"time"
)

const BaseUrl = "https://ksef-test.mf.gov.pl/api/v2"

type KSeFClient struct {
}

const GetAuthEndpoint = "/auth/challenge"

type AuthChallengeResult struct {
	Challenge string
	Timestamp time.Time
}

func (c *KSeFClient) GetAuthChallenge() (*AuthChallengeResult, error) {
	req, err := http.NewRequest("POST", BaseUrl+GetAuthEndpoint, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result AuthChallengeResult
	err = json.NewDecoder(resp.Body).Decode(&result)
	return &result, err
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

// ContextIdentifier contains the NIP
type ContextIdentifier struct {
	Nip string `xml:"Nip"`
}

func BuildAuthTokenRequestFromChallenge(challenge *AuthChallengeResult, nip string) AuthTokenRequest {
	return AuthTokenRequest{
		XmlnsXsi:  "http://www.w3.org/2001/XMLSchema-instance",
		XmlnsXsd:  "http://www.w3.org/2001/XMLSchema",
		Challenge: challenge.Challenge,
		ContextIdentifier: ContextIdentifier{
			Nip: nip,
		},
		SubjectIdentifierType: "certificateSubject",
	}
}
