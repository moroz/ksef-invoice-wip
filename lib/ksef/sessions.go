package ksef

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"time"
)

type InteractiveSession struct {
	encryptionKey   []byte
	iv              []byte
	referenceNumber *string
	validUntil      *time.Time
}

func generateRandom(len int) []byte {
	var bytes = make([]byte, len)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	return bytes
}

const EncryptionKeyLength = 256 / 8
const IVLength = 128 / 8

func (c *Client) getPublicKeyFromCertFile(pemContents []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pemContents)
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	key, ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("GetPublicKeyFromCertFile: failed to cast as RSA public key: want *rsa.PublicKey, got: %T", cert.PublicKey)
	}

	return key, nil
}

func (c *Client) GetKsefEncryptionPublicKey() (*rsa.PublicKey, error) {
	return c.getPublicKeyFromCertFile(c.EncryptionCert)
}

func encryptKeyWithRSAPublicKey(encryptionKey []byte, pubKey *rsa.PublicKey) ([]byte, error) {
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, pubKey, encryptionKey, nil)
}

func (c *Client) OpenInteractiveSession() (*InteractiveSession, error) {
	key := generateRandom(EncryptionKeyLength)
	iv := generateRandom(IVLength)

	pubKeyCert, err := c.GetKsefEncryptionPublicKey()
	if err != nil {
		return nil, err
	}

	encryptedKey, err := encryptKeyWithRSAPublicKey(key, pubKeyCert)
	if err != nil {
		return nil, err
	}

	body, err := toJSONPayload(&OpenInteractiveSessionRequest{
		FormCode: OpenInteractiveSessionRequestFormCode{
			SystemCode:    "FA (3)",
			SchemaVersion: "1-0E",
			Value:         "FA",
		},
		Encryption: OpenInteractiveSessionRequestEncryption{
			EncryptedSymmetricKey: encryptedKey,
			InitializationVector:  iv,
		},
	})
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest("POST", c.BaseUrl+"/sessions/online", body)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+*c.accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("OpenInteractiveSession: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		return nil, fmt.Errorf("OpenInteractiveSession: unexpected status code in response: want 201, got %v", resp.StatusCode)
	}

	var result OpenInteractiveSessionResponse
	err = json.NewDecoder(resp.Body).Decode(&result)

	if err != nil {
		return nil, fmt.Errorf("OpenInteractiveSession: %w", err)
	}

	return &InteractiveSession{
		encryptionKey:   key,
		iv:              iv,
		referenceNumber: &result.ReferenceNumber,
		validUntil:      &result.ValidUntil,
	}, nil
}
