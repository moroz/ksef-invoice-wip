package ksef

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io"
	"ksef-go/lib/sign"
	"log"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	*EnvironmentConfig
	env                   Environment
	certDer               []byte
	privKey               *ecdsa.PrivateKey
	nip                   string
	accessToken           *string
	refreshToken          *string
	accessTokenValidUntil *time.Time
}

func NewClient(env Environment, nip string, certDer []byte, key *ecdsa.PrivateKey) (*Client, error) {
	client := &Client{certDer: certDer, privKey: key, nip: nip, env: env}

	config, ok := EnvironmentConfigs[env]
	if !ok {
		return nil, fmt.Errorf("unknown environment %v", env)
	}

	client.EnvironmentConfig = &config

	return client, nil
}

func (c *Client) Authenticated() bool {
	return c.accessToken != nil
}

func (c *Client) Authenticate() (err error) {
	challenge, err := c.GetAuthChallenge()
	if err != nil {
		return
	}

	authTokenRequest := BuildAuthTokenRequestFromChallenge(challenge, c.nip)
	signed, err := sign.MarshalAndSign(authTokenRequest, c.privKey, c.certDer)
	if err != nil {
		return
	}
	fmt.Println(string(signed))

	xadesResp, err := c.SubmitAuthXAdESSignature(signed)
	if err != nil {
		return
	}

	err = c.AwaitAuthorization(xadesResp.ReferenceNumber, xadesResp.AuthenticationToken.Token)
	if err != nil {
		return
	}

	tokenResp, err := c.RedeemAuthToken(xadesResp.AuthenticationToken.Token)
	if err != nil {
		return
	}

	c.accessToken = &tokenResp.AccessToken.Token
	c.refreshToken = &tokenResp.RefreshToken.Token
	c.accessTokenValidUntil = &tokenResp.AccessToken.ValidUntil
	return
}

func (c *Client) GetAuthChallenge() (*AuthChallengeResult, error) {
	req, err := http.NewRequest("POST", c.BaseUrl+GetAuthEndpoint, nil)
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

func (c *Client) CheckAuthenticationStatus(ref, token string) (*CheckAuthenticationStatusResponse, error) {
	url := c.BaseUrl + strings.ReplaceAll(CheckAuthenticationStatusEndpoint, "{ref}", ref)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("CheckAuthenticationStatus: unexpected status code (want %v, got %v)", 200, resp.StatusCode)
	}

	var result CheckAuthenticationStatusResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	return &result, err
}

func (c *Client) AwaitAuthorization(ref, token string) (err error) {
	for i := range 10 {
		log.Printf("Awaiting authorization (attempt %v)", i)
		resp, err := c.CheckAuthenticationStatus(ref, token)
		if err == nil && resp.Status.Code == 200 {
			return nil
		}
		time.Sleep(time.Second)
	}
	return
}

func (c *Client) RedeemAuthToken(bearer string) (*RedeemAuthTokenResponse, error) {
	req, err := http.NewRequest("POST", c.BaseUrl+RedeemAuthTokenEndpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+bearer)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("RedeemAuthToken: unexpected status code (want %v, got %v)", 200, resp.StatusCode)
	}

	var result RedeemAuthTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	return &result, err
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

func (c *Client) SubmitAuthXAdESSignature(xmlBytes []byte) (*SubmitAuthXAdESSignatureResponse, error) {
	body := bytes.NewBuffer(xmlBytes)
	req, err := http.NewRequest("POST", c.BaseUrl+SubmitAuthXAdESSignatureEndpoint, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/xml")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 202 {
		body, _ := io.ReadAll(resp.Body)
		log.Println(string(body))

		return nil, fmt.Errorf("SubmitAuthXAdESSignature: unexpected status code (want %v, got %v)", 202, resp.StatusCode)
	}

	var result SubmitAuthXAdESSignatureResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	return &result, err
}
