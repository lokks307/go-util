package pfx

import (
	"crypto/x509"
	b64 "encoding/base64"
	"strings"

	pkcs12 "software.sslmate.com/src/go-pkcs12"
)

func DecodePFXB64(pfxDataB64 string, password string) (privateKey interface{}, certificate *x509.Certificate, err error) {

	input := strings.ReplaceAll(pfxDataB64, "\n", "")
	input = strings.ReplaceAll(input, "\r", "")

	var pfxBytes []byte

	pfxBytes, decodeErr := b64.StdEncoding.DecodeString(input)
	if decodeErr != nil {
		return nil, nil, decodeErr
	}

	privKey, cert, err := pkcs12.Decode(pfxBytes, password)
	if err != nil {
		return nil, nil, err
	}

	return privKey, cert, nil
}
