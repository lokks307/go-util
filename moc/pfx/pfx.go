package pfx

import (
	"crypto/x509"
	b64 "encoding/base64"
	"strings"

	pkcs12 "software.sslmate.com/src/go-pkcs12"
)

func DecodePFX(der []byte, password string) (privateKey interface{}, certificate *x509.Certificate, err error) {
	privKey, cert, err := pkcs12.Decode(der, password)
	if err != nil {
		return nil, nil, err
	}

	return privKey, cert, nil
}

// decode base64 encoded pfx data. This function eliminates newline and carriage return from base64 string
func DecodePFXB64(pfxDataB64 string, password string) (privateKey interface{}, certificate *x509.Certificate, err error) {

	input := strings.ReplaceAll(pfxDataB64, "\n", "")
	input = strings.ReplaceAll(input, "\r", "")

	var pfxBytes []byte

	pfxBytes, decodeErr := b64.StdEncoding.DecodeString(input)
	if decodeErr != nil {
		return nil, nil, decodeErr
	}

	return DecodePFX(pfxBytes, password)
}
