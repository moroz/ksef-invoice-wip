package ksef

import (
	"crypto/rand"
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

func (c *Client) OpenInteractiveSession() (*InteractiveSession, error) {
	key := generateRandom(EncryptionKeyLength)
	iv := generateRandom(IVLength)
}
