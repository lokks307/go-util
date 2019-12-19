package moc

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"strings"
)

func DecodePEM(pemData string) (*pem.Block, error) {
	r := strings.NewReader(pemData)

	pemBytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	data, _ := pem.Decode(pemBytes)
	if data == nil {
		return nil, errors.New("Wrong PEM format")
	}

	return data, nil
}

func ParsePemToDer(pemStr string, pswd ...string) ([]byte, error) {
	pemBlock, _ := DecodePEM(pemStr)
	if pemBlock == nil {
		return nil, errors.New("It is not PEM format")
	}

	var result []byte
	var decryptErr error
	decryptErr = nil

	if x509.IsEncryptedPEMBlock(pemBlock) { // pem could be encrypted
		if len(pswd) < 1 {
			return nil, errors.New("Password is required")
		}
		result, decryptErr = x509.DecryptPEMBlock(pemBlock, []byte(pswd[0]))
		if decryptErr != nil {
			return nil, decryptErr
		}
	} else {
		result = pemBlock.Bytes
	}

	return result, nil
}

// Parse certificate formatted in PEM or DER
func ParseDataToDer(dataB64 string, pswd ...string) []byte {
	if len(pswd) > 1 {
		return nil
	}

	var resDer []byte
	var err error

	if len(pswd) == 1 {
		resDer, err = ParsePemToDer(dataB64, pswd[0])
	} else {
		resDer, err = ParsePemToDer(dataB64)
	}
	

	if err == nil { // ... pem이다!
		return resDer
	}

	// It is not PEM format. Probably base encoded DER format
	// TODO: How to check DER is encrypted? Is decryption necessary?
	dataB64 = strings.ReplaceAll(dataB64, "\n", "")
	dataB64 = strings.ReplaceAll(dataB64, "\r", "")

	derBytes, decodeErr := base64.StdEncoding.DecodeString(dataB64)
	if decodeErr != nil { // cannot be decoded as base64, it could be corrupted data
		return nil
	}

	return derBytes
}
