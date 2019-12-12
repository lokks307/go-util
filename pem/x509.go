package pem

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

func ParseX509Cert(pemData string) (*x509.Certificate, error) {
	data, decodeErr := DecodePEM(pemData)

	if decodeErr != nil {
		return nil, decodeErr
	}

	var cert *x509.Certificate
	cert, parseErr := x509.ParseCertificate(data.Bytes)

	if parseErr != nil {
		return nil, parseErr
	}

	return cert, nil
}

func ParsePEM(dataB64 string, pswd ...string) []byte {

	if len(pswd) > 1 {
		// ??? pswd가 여러개여??
		return nil
	}

	r := strings.NewReader(dataB64)
	pemBytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil
	}

	pemBlock, _ := pem.Decode(pemBytes)

	if pemBlock != nil { // data is parsed as pem.Block. It is PEM format
		var result []byte
		var decryptErr error

		if x509.IsEncryptedPEMBlock(pemBlock) { // pem could be encrypted
			// TODO: pswd length...
			result, decryptErr = x509.DecryptPEMBlock(pemBlock, []byte(pswd[0]))
			if decryptErr != nil {
				result = nil
			}
		} else {
			result = pemBlock.Bytes
		}
		return result
	}

	// It is not PEM format. Probably base encoded DER format
	// TODO: How to check DER is encrypted? Is decryption necessary?
	dataB64 = strings.ReplaceAll(dataB64, "\n", "")
	dataB64 = strings.ReplaceAll(dataB64, "\r", "")

	pfxBytes, decodeErr := base64.StdEncoding.DecodeString(dataB64)

	if decodeErr != nil { // cannot be decoded as base64, it could be corrupted data
		return nil
	}

	return pfxBytes
}

func GetCertificateB64(dataB64 string) (*x509.Certificate, error) {
	data := ParsePEM(dataB64)

	if data == nil {
		return nil, errors.New("Can't decode data")
	}

	var cert *x509.Certificate
	cert, parseErr := x509.ParseCertificate(data)

	if parseErr != nil {
		return nil, parseErr
	}

	return cert, nil
}

//TODO: this function is temporarily located in here to use ParsePEM.
// It can be moved to ecdsa pkg or integrated into new pkg.
// If you modify this function or whole file, please inform to tethys team
func GetPrivateKey(dataB64, password string) (key interface{}, err error) {
	derBytes := ParsePEM(dataB64, password)
	privKey, error := x509.ParsePKCS8PrivateKey(derBytes)
	return privKey, error
}

func VerifyCert(pemData, CApemData string) bool {
	CAcert, parseErr := ParseX509Cert(CApemData)
	if parseErr != nil {
		return false
	}

	cert, parseErr := ParseX509Cert(pemData)
	if parseErr != nil {
		return false
	}

	verifyErr := cert.CheckSignatureFrom(CAcert)
	if verifyErr != nil {
		return false
	}

	return true
}
