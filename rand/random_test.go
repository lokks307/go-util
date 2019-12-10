package rand

import (
	"encoding/base64"
	"testing"

	"github.com/btcsuite/btcutil/base58"
)

func Test_GenRandom(t *testing.T) {
	t.Run("Base64 string", func(t *testing.T) {
		length := 32
		b64str, _ := GenRandomB64Str(length)
		decoded, _ := base64.StdEncoding.DecodeString(b64str)

		if len(decoded) != length {
			t.Fatalf("Check size - expect: %v, actual: %v", length, len(decoded))
		}
	})

	t.Run("Base58 string", func(t *testing.T) {
		length := 32
		b58str, _ := GenRandomB58Str(length)
		decoded := base58.Decode(b58str)

		if len(decoded) != length {
			t.Fatalf("Check size - expect: %v, actual: %v", length, len(decoded))
		}
	})

	t.Run("Random Seed", func(t *testing.T) {
		testCertPem := `-----BEGIN CERTIFICATE-----
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

		length := 32
		randomSeed := GenRandomSeed(testCertPem)

		if len(randomSeed) != length {
			t.Fatalf("Check size - expect: %v, actual: %v", length, len(randomSeed))
		}
	})
}
