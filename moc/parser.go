package moc

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"strings"
)

const PemBeginPhrase = "-----BEGIN "
const BeginPhraseCertificate = "-----BEGIN CERTIFICATE-----"

func GetCertificateOrPublicKey(dataB64 string) (interface{}, error) {
	if BeginPhraseCertificate == dataB64[0:27] { // this is PEM certificate
		pBlock, _ := pem.Decode([]byte(dataB64))
		if cert, err := x509.ParseCertificate(pBlock.Bytes); err == nil {
			return cert, nil
		}
		return nil, errors.New("Cannot decode PEM certificate")
	}

	dataRaw, err := base64.StdEncoding.DecodeString(dataB64)
	if err != nil { // wrong encoded data
		return nil, errors.New("Cannot decode Base64 format")
	}

	if cert, err := x509.ParseCertificate(dataRaw); err == nil {
		return cert, nil
	}

	if pubKey, err := x509.ParsePKIXPublicKey(dataRaw); err == nil {
		return pubKey, nil
	}

	return nil, errors.New("Cannot parse data to x509 or public key")
}

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
		return nil // 비번이 2개 이상?!
	}

	var resDer []byte
	var err error

	if dataB64[0:11] == PemBeginPhrase { // this file is PEM format

		if len(pswd) == 1 {
			resDer, err = ParsePemToDer(dataB64, pswd[0])
		} else {
			resDer, err = ParsePemToDer(dataB64)
		}
	} else { // It is not PEM format. Probably base encoded DER format
		// WARN: encrypted DER is not supported
		dataB64 = strings.ReplaceAll(dataB64, "\n", "")
		dataB64 = strings.ReplaceAll(dataB64, "\r", "")

		resDer, err = base64.StdEncoding.DecodeString(dataB64)
	}

	if err != nil {
		return nil
	}

	return resDer
}

func ParseHexToDer(hexStr string) []byte {
	decoded, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil
	}
	return decoded
}
