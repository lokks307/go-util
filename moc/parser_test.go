package moc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testPkcs1Key = `-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQDpO4CPuo0fXRM0U/fiNoxhLbKK0leUIJIqUnFGlj0Zn0id+19p
p/naiC33SHsffFomTnKp5RS4HIwhcZUC3io09sNGb8af0UZpFKKkblmN6Zl921Gc
Is2Ttcqfq8ZpfeVT4kC+1NqhsfuTiTzBth5+mFSh7gkC/IS/y2BNwMksKQIDAQAB
AoGAdSQhig7UWnBQ+PNDiSaEkzp0diz3N7q0LvWBV0aWIYxS2KYsYdCwOQY0sAKD
dGjsjljQVmxsX5xW9WUqxmn3H/VfUFpN6ucnbAQgb2fgiQJuBV06fWnPLiC1iNcs
oP4+XnexEE9mv1QkXSe9MqZ+3JLwwx4JlFyXEY4FevJlOAECQQD5Srggx9fM3lhW
CmXQjOgkTfaYuLWdVjXdZNX+/tLbVTkFufTOP1GVMz6hA9/3ELTMb9ndnCBNyg50
k5IqcsnBAkEA74InwlTiDCpxp3IbpZu8jQaHc+A2yBEsIhk/E52+roOP8QgEmfLV
dNMlznnxvSzvpIefqYpzsrc5D7Dh2ldsaQJAKOkFLIP/OySl9IDCUqY9FnAg7tEp
JMfYERwSLkWdTtc+g10P+qTTe5usRHpBT+dS9FXKuB5+AqYNFz58dzDdQQJBAMEh
1EvJROqIg1OCGTcm6RlTTYCsKaCU8GCiuRFpX2y8HCNB0uwNPAFTjqX8AWoJxUiI
MMB3K6rOQo4blVQKsRkCQFGggMAfAFL4gI+q/yuN3e3SQ2pGdCE/suaAAKLKgovt
kGR7Ek9bU5eBywUA0eqUI1yLOxmfdIc781ImQpGwwUY=
-----END RSA PRIVATE KEY-----`
	testPkcs1Der = `MIICXAIBAAKBgQDpO4CPuo0fXRM0U/fiNoxhLbKK0leUIJIqUnFGlj0Zn0id+19p
p/naiC33SHsffFomTnKp5RS4HIwhcZUC3io09sNGb8af0UZpFKKkblmN6Zl921Gc
Is2Ttcqfq8ZpfeVT4kC+1NqhsfuTiTzBth5+mFSh7gkC/IS/y2BNwMksKQIDAQAB
AoGAdSQhig7UWnBQ+PNDiSaEkzp0diz3N7q0LvWBV0aWIYxS2KYsYdCwOQY0sAKD
dGjsjljQVmxsX5xW9WUqxmn3H/VfUFpN6ucnbAQgb2fgiQJuBV06fWnPLiC1iNcs
oP4+XnexEE9mv1QkXSe9MqZ+3JLwwx4JlFyXEY4FevJlOAECQQD5Srggx9fM3lhW
CmXQjOgkTfaYuLWdVjXdZNX+/tLbVTkFufTOP1GVMz6hA9/3ELTMb9ndnCBNyg50
k5IqcsnBAkEA74InwlTiDCpxp3IbpZu8jQaHc+A2yBEsIhk/E52+roOP8QgEmfLV
dNMlznnxvSzvpIefqYpzsrc5D7Dh2ldsaQJAKOkFLIP/OySl9IDCUqY9FnAg7tEp
JMfYERwSLkWdTtc+g10P+qTTe5usRHpBT+dS9FXKuB5+AqYNFz58dzDdQQJBAMEh
1EvJROqIg1OCGTcm6RlTTYCsKaCU8GCiuRFpX2y8HCNB0uwNPAFTjqX8AWoJxUiI
MMB3K6rOQo4blVQKsRkCQFGggMAfAFL4gI+q/yuN3e3SQ2pGdCE/suaAAKLKgovt
kGR7Ek9bU5eBywUA0eqUI1yLOxmfdIc781ImQpGwwUY=`
	testPkcs1Password = `-----BEGIN RSA PRIVATE KEY-----
Proc-Type: 4,ENCRYPTED
DEK-Info: AES-256-CBC,01E1ACD356DE5F96470945EF9CD63485

lO/Wb3nb90Dgv3XCoZqLJEJ0zSoRxFZI1coI51RDTOC3EAaVyI1Ow7d/hvXnZV2y
wpnI/aE0vTMcQaHEqI607JF3FheBJZu4bJa2v80Wl40Riy8k+Qy/MDyq/Rtv7ju1
xphOn6QsuTj7zZgsps9/gd14WG1SL0P6gM5v9LhQQa6TU3Ew+Jln9pmxsm5DSTvw
4KX2iR5Kkb3UsBPVYnHsDcpwqHE5/z7oyezU9MdXQEX9JQX1ZA2ExTUEjjmslL7m
xqM0vkyvg08S9m97m0bYSTEcYTSzhcmVK8TPr5c1B3w1B6xpaEw9OC8LBhMMQDZl
W5gAuXR6Oy3aQ0MLC2E7l6/Y45zQSl6rnhi4MNjfGb48V8Wmehsbt0JfCMEuwEcd
6hspA2bydolzNgzYFHy5gH9iyh6tK00j3PoRJs4gk7KbuTsTeHG0fEp2zqNqyQLR
ZyNGlrIYzINRogkkjbnu6QjNIAND8YVGc5AFUuPofCR+EIk/S3ZgwvccNyUUztIl
AuTEPl2AU04D+WOkXB3SUZODVl9IC9wecFsECOQ9bCivbwnbEWFFQQWA+W4jVRN2
0+tAOjGB8v45EBJC2A3aVBITVYrAhBXGJDYn3bxiU1yElT8/X/GD/y5c/FyZ5Hz0
/oOnzRDWZ0t2hSKn/0RNEsMb54RbL3kyBmee95M+P9/jYlnTxYGq7P6cKIQcLbIl
mM85uTvD7+biNsgfwlMgm83DTiyw5uZEugCizf/wf5uyEPW/5PGHq1cIGgbdXT6B
bcNbT82dyC5gxQ4fHKCxnt8UJudkU260UxLGPgLQneKHXWmU5rRYthxpumOQvBZW
-----END RSA PRIVATE KEY-----`
	testPassword = "password"
)

func TestParser_ParseDataToDer_Success(t *testing.T) {
	keyDer1 := ParseDataToDer(testPkcs1Key)
	assert.NotNil(t, keyDer1, "pkcs1 pem keyfile parsing failed")

	keyDer2 := ParseDataToDer(testPkcs1Der)
	assert.NotNil(t, keyDer2, "pkcs1 der keyfile parsing failed")

	assert.Equal(t, keyDer1, keyDer2, "They should be equal")
}

func TestParser_ParseEncrypted_Success(t *testing.T) {
	keyDecrypt := ParseDataToDer(testPkcs1Password, testPassword)
	assert.NotNil(t, keyDecrypt, "encrypted pkcs1 pem keyfile parsing failed")

	keyDer := ParseDataToDer(testPkcs1Der)
	assert.NotNil(t, keyDer, "pkcs1 der keyfile parsing failed")

	assert.Equal(t, keyDecrypt, keyDer, "They should be equal")
}

func TestParser_ParseEncrypted_Fail(t *testing.T) {
	keyDecrypt := ParseDataToDer(testPkcs1Password, "wrong-password")

	keyDer := ParseDataToDer(testPkcs1Der)
	assert.NotNil(t, keyDer, "pkcs1 der keyfile parsing failed")

	assert.NotEqual(t, keyDecrypt, keyDer, "They should be equal")
}
