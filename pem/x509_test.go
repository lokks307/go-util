package pem

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testCertPem = `-----BEGIN CERTIFICATE-----
MIIByDCCAW6gAwIBAgIRANe5mco0f7UWwpWksW4f/UYwCgYIKoZIzj0EAwIwKjEL
MAkGA1UEAxMCQ04xCzAJBgNVBAYTAmtyMQ4wDAYDVQQKEwVsb2trczAeFw0xOTEw
MTUwMTAwNTlaFw0yMDEwMTQwMTAwNTlaMDExEjAQBgNVBAMTCXRlc3QtbmFtZTEL
MAkGA1UEBhMCa3IxDjAMBgNVBAoTBWxva2tzMFkwEwYHKoZIzj0CAQYIKoZIzj0D
AQcDQgAEgT5mmLPtodHt1/IrVDQV9Cv4JMV5ET/wtKj2IWdZ2WOP6EzYNbW4iWHP
NQ9SE+yE3XlkRvXJ+1jGP+cDReaQGaNuMGwwIQYDVR0OBBoEGAW1qy6rCwvMH44x
ZPtzTRiaRRofM8vJEzAOBgNVHQ8BAf8EBAMCAQYwEgYDVR0TAQH/BAgwBgEB/wIB
ATAjBgNVHSMEHDAagBhotF3YEXygRwDJeHGwycYHbLIuRZDs6DUwCgYIKoZIzj0E
AwIDSAAwRQIgQKw4XEGmX/nUcivfQAShcSi5fIYXy1/U1dDW4TX71OgCIQDb2m3+
4usQnMcTn4tQXSvYjoJ4J5aLZtvI1OWDC5dQEQ==
-----END CERTIFICATE-----`

	testTrueCAcertPem = `-----BEGIN CERTIFICATE-----
MIIBwDCCAWegAwIBAgIRALv1dtWfAcPkce7sXeDIBzQwCgYIKoZIzj0EAwIwKjEL
MAkGA1UEAxMCQ04xCzAJBgNVBAYTAmtyMQ4wDAYDVQQKEwVsb2trczAeFw0xOTEw
MTQwNzQ1MDBaFw0yMDEwMTMwNzQ1MDBaMCoxCzAJBgNVBAMTAkNOMQswCQYDVQQG
EwJrcjEOMAwGA1UEChMFbG9ra3MwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAAQm
zBMecRVlF/g+XyD+MUlaHBMw0mw/jIlvHGInC4AnQm4KiQkQj8K31w05EPZ4/vQ0
Zdr1KuiQaSAGLQGqrhFNo24wbDAhBgNVHQ4EGgQYaLRd2BF8oEcAyXhxsMnGB2yy
LkWQ7Og1MA4GA1UdDwEB/wQEAwIBBjASBgNVHRMBAf8ECDAGAQH/AgEBMCMGA1Ud
IwQcMBqAGGi0XdgRfKBHAMl4cbDJxgdssi5FkOzoNTAKBggqhkjOPQQDAgNHADBE
AiAPjyq+x1cpS/quxZTyMmb5HBz2GW6FXyqm3dwRl60dpQIgcxFTyoY7P/Gc8Ptz
1PB6KTQP6yJKGsLyd5ieY59Bn9o=
-----END CERTIFICATE-----`

	testFalseCAcertPem = `-----BEGIN CERTIFICATE-----
MIIB1DCCAXmgAwIBAgIRAN/nM+ZL7GV9gEt6ivJkrh4wCgYIKoZIzj0EAwIwMzEO
MAwGA1UEAxMFZmFsc2UxCzAJBgNVBAYTAnVzMRQwEgYDVQQKDAtsb2trc19mYWxz
ZTAeFw0xOTEwMTUwMTIwNDZaFw0yMDEwMTQwMTIwNDZaMDMxDjAMBgNVBAMTBWZh
bHNlMQswCQYDVQQGEwJ1czEUMBIGA1UECgwLbG9ra3NfZmFsc2UwWTATBgcqhkjO
PQIBBggqhkjOPQMBBwNCAAS2EcilTw8qwKebfD4AJDrGKIlZbubmajme5Et3dpll
yhfqpFGdP5i2z3HEXbRnzT9J5TFPHhFYLgVSM9f2KkOTo24wbDAhBgNVHQ4EGgQY
cx3szs2PEMQXjzKq5NY3ypZih4ozDJhqMA4GA1UdDwEB/wQEAwIBBjASBgNVHRMB
Af8ECDAGAQH/AgEBMCMGA1UdIwQcMBqAGHMd7M7NjxDEF48yquTWN8qWYoeKMwyY
ajAKBggqhkjOPQQDAgNJADBGAiEAwtX7m9pskUj/Y+xPT8thR/LlPVrKWxADHR3k
GSn98xMCIQCvDHAHFPn6yJ+9u9/GMMr5vUXRAPKEgGMglDkAxzGhlg==
-----END CERTIFICATE-----`

	testCertDER = `MIIDwTCCAqkCFAFLClhfd7ogpF2ghCmasR6Zp5nrMA0GCSqGSIb3DQEBCwUAMIGc
MQswCQYDVQQGEwJLUjETMBEGA1UECAwKU29tZS1TdGF0ZTEQMA4GA1UEBwwHSW5j
aGVvbjEWMBQGA1UECgwNTG9ra3MzMDcgSW5jLjEcMBoGA1UECwwTUmVzZWFyY2gg
ZGVwYXJ0bWVudDERMA8GA1UEAwwIbG9ra3MuaW8xHTAbBgkqhkiG9w0BCQEWDmNh
dGh5QGxva2tzLmlvMB4XDTE5MTIxMjAxNDEzMloXDTI5MTIwOTAxNDEzMlowgZwx
CzAJBgNVBAYTAktSMRMwEQYDVQQIDApTb21lLVN0YXRlMRAwDgYDVQQHDAdJbmNo
ZW9uMRYwFAYDVQQKDA1Mb2trczMwNyBJbmMuMRwwGgYDVQQLDBNSZXNlYXJjaCBk
ZXBhcnRtZW50MREwDwYDVQQDDAhsb2trcy5pbzEdMBsGCSqGSIb3DQEJARYOY2F0
aHlAbG9ra3MuaW8wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCtz+uv
XFJ9H99V+JV3XVGIxo4YCS8sQpHgQhfWbFbbcrKWRdMQS9Ftu2R45UBRpF9JTuoR
FqL7XcT7ZlgahdcxLX2W7usOIl4gvRgqMjQytaam9GJ1inkSYaRAgkDYNrsbxmBO
XGcIak+D6rWa9KpFxueDespHinkKjocXEeEhYTZFvMB4lxAMdVqMo8X+Y9QOlKSo
rjhlAOKByyJY67gjXBht7OJQvs1VoA84dgncZE4HIB6xaCk/rUnd8SyOoZJbNx+7
YOsimljnaH86HYdg5EYxoanjacN/yiYdgkbjAgR/rORAOXW4UZXuE7tzgGzID5Un
4j9pD98ymZdV8y47AgMBAAEwDQYJKoZIhvcNAQELBQADggEBAE2UNkvHbYQo90xz
1t/eZcs1TfE3XlehU2tXFHFPmZTtlHtFn8D3dE/sWIjFSCLX0bo1TpRLAat80yd3
c8K4M/J9RAYqIbJa4dlrsoBO9kT1RQRvz9FYQjaGZORB+5IpofwEJ3AVikaTWGxO
vzL0+0oDz/zAalPEHouxEYqdpV++g3yGh+RWBkDKnge8XryIyfWU4rp1FkcZg7nc
6TKlAdPRPH90hOISZm8Uqi3mnihz0hxm35aBBVa3ikIoO68P5xLSelwJMwulGKpo
AJl1K2H5gcWCOhit6eK7TaGA5diBtiW6584yV21KAtBmLZvsUmYz3RKJdRtsVPd0
rwHPeG0=`

	subCountry      = "kr"
	subOrganization = "lokks"
	subCommonName   = "test-name"
)

func TestPEM_PasrseX509Cert_Sucess(t *testing.T) {
	cert, parseErr := ParseX509Cert(testCertPem)
	assert.Nil(t, parseErr, "Certificate parsing failed")

	assert.Equal(t, subCountry, cert.Subject.Country[0], "They should be equal")
	assert.Equal(t, subOrganization, cert.Subject.Organization[0], "They should be equal")
	assert.Equal(t, subCommonName, cert.Subject.CommonName, "They should be equal")
}

func TestPEM_PasrseX509Cert_Fail(t *testing.T) {
	_, parseErr := ParseX509Cert("this is not a pem format")

	assert.NotNil(t, parseErr, "This case must make error but no error")
}

func TestPEM_VerifyCert_True(t *testing.T) {
	check := VerifyCert(testCertPem, testTrueCAcertPem)

	assert.True(t, check, "Verification must succeed")
}

func TestPEM_VerifyCert_False(t *testing.T) {
	check := VerifyCert(testCertPem, testFalseCAcertPem)

	assert.False(t, check, "Verification must fail")
}

func TestPEM_GetCertificateB64_Success(t *testing.T) {
	cert, parseErr := GetCertificateB64(testCertPem)

	assert.Nil(t, parseErr, "[PEM]: Certificate parsing failed")

	assert.Equal(t, subCountry, cert.Subject.Country[0], "[PEM]: They should be equal")
	assert.Equal(t, subOrganization, cert.Subject.Organization[0], "[PEM]: They should be equal")
	assert.Equal(t, subCommonName, cert.Subject.CommonName, "[PEM]: They should be equal")

	cert, parseErr = GetCertificateB64(testCertDER)
	assert.Nil(t, parseErr, "[DER]: Certificate parsing failed")
	assert.Equal(t, "lokks.io", cert.Subject.CommonName, "[PEM]: They should be equal")
}
