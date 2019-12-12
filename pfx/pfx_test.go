package pfx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	// generated using openssl, rsa 2048bit
	rsaBase64 = `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
yrPsT7a2/wQIHommi7KIICUCAggA`
	// generated using openssl, ecdsa prime256v1 params
	ecBase64 = `MIIEOgIBAzCCBAAGCSqGSIb3DQEHAaCCA/EEggPtMIID6TCCAt8GCSqGSIb3DQEH
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
	pswd            = "password"
	certSubjectInfo = "CN=lokks.io,OU=Research department,O=Lokks307 Inc.,L=Incheon,ST=Some-State,C=KR"
)

func TestDecodeB64(t *testing.T) {
	_, cert, err := DecodePFXB64(rsaBase64, pswd)
	assert.Nil(t, err, "[RSA Test]: Decoding error")
	assert.Equal(t, cert.Subject.String(), certSubjectInfo, "[RSA Test]: Wrong subject information")

	_, cert, err = DecodePFXB64(ecBase64, pswd)
	assert.Nil(t, err, "[ECDSA Test]: Decoding error")
	assert.Equal(t, cert.Subject.String(), certSubjectInfo, "[ECDSA Test]: Wrong subject information")
}

func TestDecodeFail(t *testing.T) {
	_, _, err := DecodePFXB64(rsaBase64, "wrong_password")
	assert.NotNil(t, err, "[RSA Test]: Decoding must fail. But there is no error")

	_, _, err = DecodePFXB64(ecBase64, "wrong_password")
	assert.NotNil(t, err, "[ECDSA Test]: Decoding must fail. But there is no error")
}
