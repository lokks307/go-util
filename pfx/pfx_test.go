package pfx

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testBase64 = `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
cvSzwtFWNkyQXfE6VYu/BAhq5OgX2kr3DAICCAA=`
)

func TestDecode(t *testing.T) {
	var pswd = ""
	_, cert, err := DecodePFXB64(testBase64, pswd)

	assert.Nil(t, err, "Decoding error")

	fmt.Println(cert.Subject)
}

// func TestBase64(t *testing.T) {
// 	// 	var b64txt string = `SGVsbG8sIGdvbGFuZy
// 	// ENCkknbSBiYXNlNjQ=`
// 	// b64txt = strings.ReplaceAll(b64txt, "\n", "")
// 	// b64txt = strings.ReplaceAll(b64txt, "\r", "")

// 	block, rest := pem.Decode([]byte(testBase64))

// 	fmt.Println(block)
// 	fmt.Println(rest)

// 	// var dst []byte
// 	// len, _ := base64.StdEncoding.Decode(dst, block.Bytes)

// 	// fmt.Println(len)

// 	var pswd = ""
// 	_, cert, _ := pkcs12.Decode(block.Bytes, pswd)

// 	fmt.Println(cert)

// }
