package rand

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/pbkdf2"
)

// GenRandomBytes generate random bytes
func GenRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenRandomB64Str generate random base64 string
func GenRandomB64Str(n int) (string, error) {
	b, err := GenRandomBytes(n)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(b), nil
}

// GenRandomB58Str generate random base58 string
func GenRandomB58Str(n int) (string, error) {
	b, err := GenRandomBytes(n)
	if err != nil {
		return "", err
	}

	return base58.Encode(b), nil
}

const (
	randomSeedIter = 10000
	randomSeddLen  = 32
)

// GenRandomSeed generate random seed using pbkdf2
func GenRandomSeed(pem string) []byte {
	return pbkdf2.Key([]byte(pem), []byte(""), randomSeedIter, randomSeddLen, sha256.New)
}
