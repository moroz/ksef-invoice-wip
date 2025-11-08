package main

import (
	"net/http"
)

const BaseUrl = "https://ksef-test.mf.gov.pl/api/v2"

type KSeFClient struct {
}

const GetAuthEndpoint = "/auth/challenge"

func (c *KSeFClient) GetAuthChallenge() (*http.Response, error) {
	req, err := http.NewRequest("POST", BaseUrl+GetAuthEndpoint, nil)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req)
}
