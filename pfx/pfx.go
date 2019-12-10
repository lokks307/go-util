package pfx

import (
	"crypto/x509"
	b64 "encoding/base64"
	"regexp"

	pkcs12 "software.sslmate.com/src/go-pkcs12"
)

func DecodePFX(pfxDataB64 string, password string) (privateKey interface{}, certificate *x509.Certificate, err error) {

	re := regexp.MustCompile(`\r?\n`)
	input := re.ReplaceAllString(pfxDataB64, "")

	pfxBytes, decodeErr := b64.StdEncoding.DecodeString(input)
	if decodeErr != nil {
		panic(decodeErr)
	}

	privKey, cert, err := pkcs12.Decode(pfxBytes, password)
	if err != nil {
		return nil, nil, err
	}

	return privKey, cert, nil
}
