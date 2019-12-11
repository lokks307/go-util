package ecdsa

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testMsg = "Hello ecdsa"

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
)

func TestEcdsa_SignVerifyNoPasswordCase_Success(t *testing.T) {
	signature, err := Sign([]byte(testMsg), testSkNoPass, "")
	if err != nil {
		t.Error(err)
	}

	success := Verify([]byte(testMsg), signature, testNoPassCert)

	assert.True(t, success, "Verification must succeed.")
}

func TestEcdsa_SignVerifyNoPasswordCase_Fail_WrongPem(t *testing.T) {
	_, err := Sign([]byte(testMsg), "this is not a pem format", "")
	assert.NotNil(t, err, "This case must make error but no error")
}

// 19.12.11: lokks307/pkcs8로 변경 후 테스트 통과
func TestEcdsa_CaseOfPassword(t *testing.T) {
	signature, err := Sign([]byte(testMsg), testSkPass, "password")
	if err != nil {
		t.Error(err)
	}

	success := Verify([]byte(testMsg), signature, testPassCert)
	if !success {
		t.Errorf("TestEcdsa_CaseOfNoPassword failed: want '%t', got '%t'", true, success)
	}
}
