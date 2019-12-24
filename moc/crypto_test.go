package moc

import (
	go_ecdsa "crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	go_rsa "crypto/rsa"
	"encoding/base64"
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testEncPkcs8Key = `-----BEGIN ENCRYPTED PRIVATE KEY-----
MIIC3TBXBgkqhkiG9w0BBQ0wSjApBgkqhkiG9w0BBQwwHAQIRGOcAqyKvuQCAggA
MAwGCCqGSIb3DQIJBQAwHQYJYIZIAWUDBAEqBBDGkQGTNCUO0jkwXBKsuQXIBIIC
gD4Iq/3bYF4iYX/HDMBJbmT8n7hW7mQnaw3eA3vHXh5vg8eZGyjFpThrZ+C+OKTn
CsrEbYRntV2g/7iWmv8snedaIcH09L68f4gYOOOqFb47jqVAmjLuP/TY5HS7bYki
CKG5vCx4hN9avvFu3M27bjhCZn1aGRa+0ILMjEhVZQCNb1BXbLQq4k5a7Hgn+1VP
7l1jcyxXyj3JpezQz080J6IsYdLPcnb2g53foDm49Dqp3C1s1/V5Bu8P8+rrQeqt
RAcN58i3t3n+x+d/y2608paTHw+FQ5POc72QUHaZnDH96C1pAr3ku7m8S4uqquO6
Fwnsgm+huuD5dPsyPlsILyypnJ2cCuHtUaPaofgGQbMLyEcLpAxVi4+cbyghNeEf
cauftLV8eNpzSWVZfv0JSRkWLDtIhoPUgGPFE/hl0e6PZ0adKDhiRfv+9J7Apms7
s0TE5V2tXSg7evkspP9khgroOGn3nog238YgVRSJr6FymgjyIbIy6aPsoM/Ugl7r
VPRAcTQ3EVmqOf9wlqAX4qgPb+sLcoJTrmeTyMlxcufNMzQhZsslx0IXAUXa39Am
coX1+b1465IecnOCk2iUCPd8rvTEYUpmttQztj8gjN2me7iMzvM1niUmT8xxJ/Dl
JUalOHijOvEqQyNVTlSHv9AqIMuTOzmkzJcH2HwQ8mi3s765Lf85ufF36H5vuXMq
eEzoyKyCQGmMCfhU5RKQnL6mumYivqucgtNWbBsKZcfbBPBTAJpfz36U1F4BVsO1
6KNOmYkf3HsPP+Hj4zUywfj/xPtDVGnlmtkNnEeitpeEtSaruY1JJeXVt60D5Ijd
0RJ+DdOpDggk+XxWtR91Tpg=
-----END ENCRYPTED PRIVATE KEY-----`
	testPfxECBase64 = `MIIEOgIBAzCCBAAGCSqGSIb3DQEHAaCCA/EEggPtMIID6TCCAt8GCSqGSIb3DQEH
BqCCAtAwggLMAgEAMIICxQYJKoZIhvcNAQcBMBwGCiqGSIb3DQEMAQYwDgQI8ynL
P19WDvACAggAgIICmOhHmvhFn/BJfujHzKhuAgNFC2HlLPdLxj5scBOCOp4bQOz6
32tRR8un78f/Jbb3gObOwzGvwM3jhJxeHllLCEnb6rdJMKr9l1IGhhEDwXdeHfs3
7HJ+l+ajYemd/qsw2AeQc/Oi4lYGx86u1GVN2ppOIu1kauejwMYuVYKMopm7g11Q
rRIrJWqBIX2wm4GKZ7Bc1MnYdpALSoxSIYmAd17XAv3rKhHK4bU+BmmswQrInOb+
p5R/RP2DV+OlHkAbE8JhgJtNoVzQ50O31dOR1D96LeMrwfTsQNOf4IrLp9uQPUBK
LJT3XBtdk6NOnJZYvtu/XaO4Z1qoh3LLeIglWW4hTjbDzAT+Se5OBf0nv+Xpman3
rJAxcSJgjAehHQ8V/l1og5VGLfAo967d8PKYtubaWQmfGd2W7oR6pJa6Y+0jnoI9
ZKcAREke62uQwA24ClJtKYNrEbt4MVzTQ/JOvEIf+yYAmnhNBSdkHkF5xSGXzhun
oGmrEGhwM4b8c72OG3KXqoaGppEaf2uvC4xoFFUtxROjXNSq11iCCSFcUB68so/+
+gqAGRo1Qznahs1dg8Mu+ZWP/8LC4FBDZs8xjsMKXkHgbT38aXmbaGsNBHb8Wkcf
9P8EgWjkgQ6o0A90KvxDXPfn6GLhMxQcJezVrqE8AvGr85u/cEdRDFjWnhU0zxql
b89V2LfxKln+/9e5Z6+MsYZT69b0FrgLGCZXZB7HIhPUZDH8ByrPFoapoGZLZx56
XrQ+UPTVK7b4s6+pz+0ShX1e8MVASfTSvtVG7RUm5RuirxqVmYHrJ80YlV7fR/wh
kVxrvOwEFlVRhuEWWLoj5Ke8lUlHVLuaX6VsmUAuOSvW7aSn7gdnvfBwG1cJfv4F
esojgPkwggECBgkqhkiG9w0BBwGggfQEgfEwge4wgesGCyqGSIb3DQEMCgECoIG0
MIGxMBwGCiqGSIb3DQEMAQMwDgQI/Il7zB7o/8kCAggABIGQyHiuDqhGC0ot75Fi
TM6fXiLJLTmuLWL7dMWtJMqV+iauAyccgGW9GriP99dL/W7AkY8DixEEbjCDIe6g
KiYGFsMKdV9j04HCea2UGnVac10U28eOFJjwGw1nk9drupTbWfBVviuwo0TCDtb1
1ZPmryIq5gVyBzl7WYgkXloEO/sXQgcYHrWM3QEYPtmxBmwsMSUwIwYJKoZIhvcN
AQkVMRYEFLNZPMcIvGUa+bnmDw6MwTDgY5rYMDEwITAJBgUrDgMCGgUABBRH+UYs
rzoNucFaQqM7O7wFCv0oowQIpkpdLledyvgCAggA`
	testSkNoPass = `-----BEGIN PRIVATE KEY-----
MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQgJOMyeyvaMej1Ar2G
TlRmcvmqH2dwrpR3ua8Nh0n/PYehRANCAARRvKeSWmyEVpMqz0Qop9ElZ0nkdV0T
IcAbuiwD2BpTZ6JSPOTKMB6c3LN88rdSaDvt0uDjxpt20NIT6+zQyQmI
-----END PRIVATE KEY-----`

	testNoPassCert = `-----BEGIN CERTIFICATE-----
MIIByDCCAW6gAwIBAgIRAMzigilbK+Bbi4PSb6nIfSQwCgYIKoZIzj0EAwIwKjEL
MAkGA1UEAxMCQ04xCzAJBgNVBAYTAmtyMQ4wDAYDVQQKEwVsb2trczAeFw0xOTEy
MDMwNzQzMTdaFw0yMDEyMDIwNzQzMTdaMDExEjAQBgNVBAMTCXRlc3QtbmFtZTEL
MAkGA1UEBhMCa3IxDjAMBgNVBAoTBWxva2tzMFkwEwYHKoZIzj0CAQYIKoZIzj0D
AQcDQgAEUbynklpshFaTKs9EKKfRJWdJ5HVdEyHAG7osA9gaU2eiUjzkyjAenNyz
fPK3Umg77dLg48abdtDSE+vs0MkJiKNuMGwwIQYDVR0OBBoEGG/FFsvAqPRaRmMq
kxU3MhmvzKoug9Py6zAOBgNVHQ8BAf8EBAMCAQYwEgYDVR0TAQH/BAgwBgEB/wIB
ATAjBgNVHSMEHDAagBhotF3YEXygRwDJeHGwycYHbLIuRZDs6DUwCgYIKoZIzj0E
AwIDSAAwRQIgEijSkDoAXD+8WdBrK3i0tCSZpq+JcCF9N9Rxeqn8AnMCIQC8xh3i
W1LMjET4Ca0Dhmtdmj0ziUCAPDh6q/p1ImRxWg==
-----END CERTIFICATE-----`

	testSkPass = `-----BEGIN ENCRYPTED PRIVATE KEY-----
MIH0MF8GCSqGSIb3DQEFDTBSMDEGCSqGSIb3DQEFDDAkBAyZjG7Iq7gZUk8UqRUC
AwX7QAIBEDAMBggqhkiG9w0CCQUAMB0GCWCGSAFlAwQBAgQQKQjbYz2Tdj2IAiMY
j1TnYgSBkB+XjoEBpfzDS3tWaxj5Rewe7vslpdH95WBx6a2m4XTAQDLMrBjEPd0T
Nmaohe+vusI8bVWkmyTU4xSWdUzqbxf65+CsuBCnPgr/zKtaiQvqpAdmLq64OnM0
s8UAFpPi/kUtvthsizEI6lPVIGnLGQHMLKFvtvlHX4bcdsZdYHKXZv4sZpIwloz8
GS9JbsVawA==
-----END ENCRYPTED PRIVATE KEY-----`

	testPassCert = `-----BEGIN CERTIFICATE-----
MIIBxzCCAW6gAwIBAgIRAMzKzV/buR4DAIAICWlMZaYwCgYIKoZIzj0EAwIwKjEL
MAkGA1UEAxMCQ04xCzAJBgNVBAYTAmtyMQ4wDAYDVQQKEwVsb2trczAeFw0xOTEy
MDQwMzQxNTFaFw0yMDEyMDMwMzQxNTFaMDExEjAQBgNVBAMTCXRlc3QtbmFtZTEL
MAkGA1UEBhMCa3IxDjAMBgNVBAoTBWxva2tzMFkwEwYHKoZIzj0CAQYIKoZIzj0D
AQcDQgAEiH0p4x20X1BP/FOJXwl/TRItrZl8Vawx/9hoTQnAu9plWteedHKLxZCP
zBMFoCaxlUHP+oZSMMtXGuItwU7hEqNuMGwwIQYDVR0OBBoEGHQGat9d3Yg99SWS
LDCOrEi0NiTBMPn+zTAOBgNVHQ8BAf8EBAMCAQYwEgYDVR0TAQH/BAgwBgEB/wIB
ATAjBgNVHSMEHDAagBhotF3YEXygRwDJeHGwycYHbLIuRZDs6DUwCgYIKoZIzj0E
AwIDRwAwRAIgZCNKWjire6lkJvQvmOhpKwM9fJZn5ViJZEQyRP2q6MMCIAGtHn2j
69zpemNZ7LusM3bqgB4gt+0kabQgAAHsoSDd
-----END CERTIFICATE-----`

	testEddsaKey = `-----BEGIN PRIVATE KEY-----
MC4CAQAwBQYDK2VwBCIEIE9Lop/DzZRuESo5HwOnDRJg9vbkA3Rvxmlj6nslg5ed
-----END PRIVATE KEY-----`

	testEddsaCert = `-----BEGIN CERTIFICATE-----
MIIBoDCCAVICAQEwBQYDK2VwMHwxCzAJBgNVBAYTAktSMRAwDgYDVQQIDAdJbmNo
ZW9uMRMwEQYDVQQHDApOYW1kb25nLWd1MRUwEwYDVQQKDAxMb2trczMwNyBJbmMx
HDAaBgNVBAsME1Jlc2VhcmNoIERlcGFydG1lbnQxETAPBgNVBAMMCGxva2tzLmlv
MB4XDTE5MTIyNDAyMjcwMloXDTIwMTIyMzAyMjcwMlowfDELMAkGA1UEBhMCS1Ix
EDAOBgNVBAgMB0luY2hlb24xEzARBgNVBAcMCk5hbWRvbmctZ3UxFTATBgNVBAoM
DExva2tzMzA3IEluYzEcMBoGA1UECwwTUmVzZWFyY2ggRGVwYXJ0bWVudDERMA8G
A1UEAwwIbG9ra3MuaW8wKjAFBgMrZXADIQD0kOHvErHEs8TWNdDcSRMrhcQ4e030
OBQTq25a7o4dRTAFBgMrZXADQQDPwoEv9CjuL/OsQBi4lww0E6xnMz+kowcxFBU+
P8Csa5gLeiD2d9TY2oaNXW2QgfIULZ2mrwHtM+RIioJ39GcE
-----END CERTIFICATE-----`

	testPfxCertSubject = "CN=lokks.io,OU=Research department,O=Lokks307 Inc.,L=Incheon,ST=Some-State,C=KR"

	testMsg = "Hello ecdsa"
)

func TestCrypto_PasrseRSAKey_Sucess(t *testing.T) {
	rsaKey, err := GetPrivateKey(testPkcs1Key, "")
	assert.Nil(t, err, "pkcs1 pem keyfile parsing failed")
	assert.IsType(t, &(rsa.PrivateKey{}), rsaKey, "[pem]: type should be *rsa.PrivateKey")

	rsaKey2, err := GetPrivateKey(testPkcs1Der, "")
	assert.Nil(t, err, "pkcs1 der keyfile parsing failed")
	assert.IsType(t, &(rsa.PrivateKey{}), rsaKey2, "[der]: type should be *rsa.PrivateKey")

	assert.Equal(t, rsaKey, rsaKey2, "rsa key must be equal")

	rsaKey3, err := GetPrivateKey(testPkcs1Password, testPassword)
	assert.Nil(t, err, "encrypted pkcs1 pem keyfile parsing failed")
	assert.IsType(t, &(rsa.PrivateKey{}), rsaKey3, "[enc-pem]: type should be *rsa.PrivateKey")

	assert.Equal(t, rsaKey3, rsaKey, "rsa key must be equal")
}

func TestCrypto_ParsePKCS8RSA_Success(t *testing.T) {
	rsaKey, err := GetPrivateKey(testPkcs1Key, "")
	assert.Nil(t, err, "pkcs1 pem keyfile parsing failed")
	assert.IsType(t, &(go_rsa.PrivateKey{}), rsaKey, "[pem]: type should be *rsa.PrivateKey")

	rsaKey2, err := GetPrivateKey(testEncPkcs8Key, testPassword)
	assert.Nil(t, err, "pkcs8 pem keyfile parsing failed")
	assert.IsType(t, &(go_rsa.PrivateKey{}), rsaKey2, "[pem]: type should be *rsa.PrivateKey")

	assert.Equal(t, rsaKey2, rsaKey, "rsa key must be equal")
}

func TestCrypto_ParsePFX_Success(t *testing.T) {
	key, cert, err := DecodePFXB64(testPfxECBase64, testPassword)
	assert.Nil(t, err, "[PFX Test]: Decoding error")
	assert.Equal(t, cert.Subject.String(), testPfxCertSubject, "[PFX Test]: Wrong subject information")
	assert.IsType(t, &(go_ecdsa.PrivateKey{}), key, "[PFX Test]: type should be *ecdsa.PrivateKey")
}

func TestCrypto_SignVerifyNoPasswordCase_Success(t *testing.T) {
	signature, err := Sign([]byte(testMsg), testSkNoPass, "")
	assert.Nil(t, err, "Signature generation failed")

	success := Verify([]byte(testMsg), signature, testNoPassCert)
	assert.True(t, success, "Verification must succeed.")
}

func TestCrypto_SignVerifyNoPasswordCase_Fail_WrongPem(t *testing.T) {
	_, err := Sign([]byte(testMsg), "this is not a pem format", "")
	assert.NotNil(t, err, "This case must make error but no error")
}

// 19.12.11: lokks307/pkcs8로 변경 후 테스트 통과 확인
func TestCrypto_CaseOfPassword(t *testing.T) {
	signature, err := Sign([]byte(testMsg), testSkPass, testPassword)
	assert.Nil(t, err, "Signature generation failed")

	success := Verify([]byte(testMsg), signature, testPassCert)
	assert.True(t, success, "Verification must succeed.")
}

func TestCrypto_EddsaSignature(t *testing.T) {
	key, err := GetPrivateKey(testEddsaKey, "")
	assert.Nil(t, err, "Private key parsing failed")
	assert.IsType(t, ed25519.PrivateKey{}, key, "Type must be ed25519.PrivateKey")

	signature, err := Sign([]byte(testMsg), testEddsaKey, "")
	assert.Nil(t, err, "Signature generation failed")

	fmt.Println(base64.StdEncoding.EncodeToString(signature))

	success := Verify([]byte(testMsg), signature, testEddsaCert)
	assert.True(t, success, "Verification must succeed.")
}

type ECDSA struct {
	RBigInt *big.Int
	SBigInt *big.Int
}

func TestCrypto_Asn1(t *testing.T) {
	var ecdsaDerSigB64 = "MEYCIQCbciuwShp4Sm8RLQQXwziLDM5yYY/H75cp9V6O0AZSBAIhAPuZXfxRGJSWZzN5kfV0QnZ7/Qn34UKwZJyRv/qVmmrW"
	sigRaw, parseErr := base64.StdEncoding.DecodeString(ecdsaDerSigB64)
	assert.Nil(t, parseErr, "Signature parsing failed")

	success := Verify([]byte(testMsg), sigRaw, testNoPassCert)
	assert.True(t, success, "Verification must succeed.")
}
