package ecdsa

import (
	"crypto/ecdsa"
	"crypto/rand"
	"math/big"

	"github.com/lokks307/go-util/pem"

	"github.com/youmark/pkcs8"
)

// TODO: write log

func Sign(msg []byte, skPem, pass string) ([]byte, error) {
	var privateKey *ecdsa.PrivateKey
	var err error

	privateKey, err = GetECPrivateKeyFromPem(skPem, pass)

	if err != nil {
		return nil, err
	}

	r, s, err := ecdsa.Sign(rand.Reader, privateKey, msg)
	if err != nil {
		return nil, err
	}

	signature := r.Bytes()
	signature = append(signature, s.Bytes()...)

	return signature, nil
}

func Verify(msg, sigBytes []byte, certPem string) bool {
	publicKey, err := GetECPublicKeyFromPem(certPem)

	if err != nil {
		return false
	}
	halfSigLen := len(sigBytes) / 2

	r := new(big.Int)
	r.SetBytes(sigBytes[:halfSigLen])

	s := new(big.Int)
	s.SetBytes(sigBytes[halfSigLen:])

	return ecdsa.Verify(publicKey, msg, r, s)
}

func GetECPublicKeyFromPem(pemData string) (*ecdsa.PublicKey, error) {
	cert, parseErr := pem.ParseX509Cert(pemData)

	if parseErr != nil {
		return nil, parseErr
	}

	publicKey := cert.PublicKey.(*ecdsa.PublicKey)

	return publicKey, nil
}

func GetECPrivateKeyFromPem(pemData, pass string) (*ecdsa.PrivateKey, error) {
	data, decodeErr := pem.DecodePEM(pemData)

	if decodeErr != nil {
		return nil, decodeErr
	}

	var err error
	var ecPrivateKey *ecdsa.PrivateKey

	if pass != "" {
		ecPrivateKey, err = pkcs8.ParsePKCS8PrivateKeyECDSA(data.Bytes, []byte(pass))
	} else {
		ecPrivateKey, err = pkcs8.ParsePKCS8PrivateKeyECDSA(data.Bytes)
	}

	if err != nil {
		return nil, err
	}

	return ecPrivateKey, nil
}
