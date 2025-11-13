package config

import (
	"log"
	"os"
)

func MustGetenv(name string) string {
	val := os.Getenv(name)
	if val == "" {
		log.Fatalf("Environment variable %s is not set!", name)
	}
	return val
}

var NipNumber = MustGetenv("NIP_NUMBER")
var AuthenticationCertPath = MustGetenv("AUTHENTICATION_CERT_PATH")
var AuthenticationPrivKeyPath = MustGetenv("AUTHENTICATION_PRIV_KEY_PATH")
var AuthenticationPrivKeyPassphrase = MustGetenv("AUTHENTICATION_PRIV_KEY_PASSPHRASE")
