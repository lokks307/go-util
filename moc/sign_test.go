package moc

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/btcsuite/btcutil/base58"
	"github.com/lokks307/go-util/bytesbuilder"
	"github.com/stretchr/testify/assert"
)

const (
	WorldName          = "LOKKS307ENTERPRISE"
	WorldTime          = "1581984000"
	WorldAnonymousUser = true
	WorldKeyCcyName    = "TYS"
	WorldKeyCcyAmount  = "100000000"
	WorldAuthorityId   = "*"
	WorldEdenChainId   = "RY89zlPUvZsNJowp5P20BloXrIxWoJLcQwoti3lmNeU="
	WorldCreatorId     = "3k1iyMAUFpE8JTmWVBZk6hHw7CyVWm6VWFYTwBZDMPB1"

	ChainName           = "DOSEEASE@LOKKS307"
	ChainTime           = "1581984000"
	ChainCustomContract = false
	ChainOracle         = false
	ChainDym            = false
	ChainScript         = false

	CreatorCN   = "(주)록스307(48935)000302620190801162754379"
	CreatorCert = `-----BEGIN CERTIFICATE-----
MIIFuzCCBKOgAwIBAgIEJXcUoDANBgkqhkiG9w0BAQsFADBSMQswCQYDVQQGEwJr
cjEQMA4GA1UECgwHeWVzc2lnbjEVMBMGA1UECwwMQWNjcmVkaXRlZENBMRowGAYD
VQQDDBF5ZXNzaWduQ0EgQ2xhc3MgMjAeFw0xOTA3MzExNTAwMDBaFw0yMDA4MDEx
NDU5NTlaMIGLMQswCQYDVQQGEwJrcjEQMA4GA1UECgwHeWVzc2lnbjEUMBIGA1UE
CwwLY29ycG9yYXRpb24xDDAKBgNVBAsMA0lCSzEOMAwGA1UECwwFNDg5MzUxNjA0
BgNVBAMMLSjso7wp66Gd7IqkMzA3KDQ4OTM1KTAwMDMwMjYyMDE5MDgwMTE2Mjc1
NDM3OTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANC7FbjMJzChoFMV
ZjmtZKetExTzuTbHThu6dJ647jux3dMHoYaoQsk8eh6Fvk4jQBYoQdrQw8rM3ghX
gCrsMo3gEuGAYG4SD8aYMVZaH2l49KveTbJ/XYe2F/Jccg6gpnVI5Ua57QyyucQI
zIvme0QBCVh2UhYLdRn02QhgQZBCykOYipbPHAZTvRTcmPRr8LhmRBeibRWwH0B/
lUdSuz0LeVoWeuFkAZdVm0ovQvDempUQ/siWaf7Pv+F5h/UECotiZfsv9+tXrImQ
eKih7O2koBITzrCtFOEl0k7VGUpgRsj9Ylal2ljcmo6AICWuUx32cnJE7PDIyPjA
sR9J+1UCAwEAAaOCAl0wggJZMIGPBgNVHSMEgYcwgYSAFO/cRNLGjcAOozjAfJPG
w0G/So/woWikZjBkMQswCQYDVQQGEwJLUjENMAsGA1UECgwES0lTQTEuMCwGA1UE
CwwlS29yZWEgQ2VydGlmaWNhdGlvbiBBdXRob3JpdHkgQ2VudHJhbDEWMBQGA1UE
AwwNS0lTQSBSb290Q0EgNIICEBwwHQYDVR0OBBYEFI0OLYOzYtbWZCjeezR2qdwK
xKmJMA4GA1UdDwEB/wQEAwIGwDB5BgNVHSABAf8EbzBtMGsGCSqDGoyaRQEBAjBe
MC4GCCsGAQUFBwICMCIeIMd0ACDHeMmdwRyylAAgrPXHeMd4yZ3BHAAgx4WyyLLk
MCwGCCsGAQUFBwIBFiBodHRwOi8vd3d3Lnllc3NpZ24ub3Iua3IvY3BzLmh0bTBt
BgNVHREEZjBkoGIGCSqDGoyaRAoBAaBVMFMMDijso7wp66Gd7IqkMzA3MEEwPwYK
KoMajJpECgEBATAxMAsGCWCGSAFlAwQCAaAiBCC5C6DsN80Bi4QJ+OaskFfdmstv
PXh8AGGBIjpNrT7j1TByBgNVHR8EazBpMGegZaBjhmFsZGFwOi8vZHMueWVzc2ln
bi5vci5rcjozODkvb3U9ZHA1cDQ4NTE5LG91PUFjY3JlZGl0ZWRDQSxvPXllc3Np
Z24sYz1rcj9jZXJ0aWZpY2F0ZVJldm9jYXRpb25MaXN0MDgGCCsGAQUFBwEBBCww
KjAoBggrBgEFBQcwAYYcaHR0cDovL29jc3AueWVzc2lnbi5vcmc6NDYxMjANBgkq
hkiG9w0BAQsFAAOCAQEAGUai0Q495+XABkdrrbdy1yvnYdM5Jn6c5vPjlFkmgxBn
ga6F24di2V+4OTt4EnxHTjUE1h/2PdMTHm2K9RBK3StCexISbUkPZZJnZtAdPAxC
cz11//C0WUKv/gMVArK9eSUCGm8xNOlIoknoo8Q7QNthVndQcXtvMF8cvitpPSzv
5TmJ3oC9uzvWsOefd4/s7dV3OXktQfeQlftBudZ7dsAdLezYy+PWFlYZjd87fJNx
f45aDHew7aB+Y26k9X1XHeJnS9ITq7Qb9vOp8WlPbKL3/btIOwK4n5eY+rxBOf+y
Kr6BqrwntHsZiM0rX4BJRqBedexNUQfnBibTX8FwwQ==
-----END CERTIFICATE-----`

	CreatorKey = `-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDQuxW4zCcwoaBT
FWY5rWSnrRMU87k2x04bunSeuO47sd3TB6GGqELJPHoehb5OI0AWKEHa0MPKzN4I
V4Aq7DKN4BLhgGBuEg/GmDFWWh9pePSr3k2yf12HthfyXHIOoKZ1SOVGue0MsrnE
CMyL5ntEAQlYdlIWC3UZ9NkIYEGQQspDmIqWzxwGU70U3Jj0a/C4ZkQXom0VsB9A
f5VHUrs9C3laFnrhZAGXVZtKL0Lw3pqVEP7Ilmn+z7/heYf1BAqLYmX7L/frV6yJ
kHiooeztpKASE86wrRThJdJO1RlKYEbI/WJWpdpY3JqOgCAlrlMd9nJyROzwyMj4
wLEfSftVAgMBAAECggEBAJSYm54zuLKmXbzIPz5QcpfsXulLmU+mE7UpyFw2y2Mz
/Pd/Nz7mCqW4qDeSfyihb75Waouck8aMkoTdxiDIhjT6kHZ5Li0uLozzTCxtfG9Z
7NVuLWIlyjATMnkX1xsSw247tv3i3x9rHVSx7uyp6JdBJaYUldubzIHuDZQo8oP/
vgDma5USMGIUPnI1kirJqUYldToTXU98t7y1xOMYgn/xeUwoc/m/ggDNZnMYx9nL
74gOvUkAUfteK6pBH6e46LVXx8sIm/tRFE1gY2oDVBamOmNUJx8+ATvcLUdrePCq
saQSInUPaWEBSN3jJ3tF0yTn7zQVYjwf94Hjg54NGwECgYEA/SXRJVN3tNZgTmFh
g/2D51wztjzhTziRfk/xJCcvDr9ix/wMWKcuVsqzuZ5E1o8len+SWZR6d9aEpaZT
HVv3qzd6MJQjp3v2iurkyjVElRFKv66Ml9o2ogPOsxT2+QaRxA2R+eNFSoyBCjxL
lpGhX3fi0P6HzFiaM7YKr75hYjcCgYEA0xUmqyt5b63+fcQNJJbiqLPoe8WjI6Rw
aSPtmV+KAgqK4X11lvpduTII6LEv+2FDbvySSAOpUVjXWsyjYIqEB6I24rAAKRU7
0nC63XMm7ELsOCC6E3BrvEq+I5iL84MsFXA/0941jWz1JQthGyUvSXMvQy+WYGkd
F7Vk2xfGONMCgYA3buOIh+mLqPH31+0xqP8MV61fN5+i6GE7xqeoYgg084XfAvt1
Ik7MZKDMgbLTRQ9Q7sSOZywWN9xeJbjHkg7+6CwSnS3djBClxOAKw3VcKygyJzfU
PM1/1tOZdXrLdzvOMaaIkNLoizHqmt92fjdXH3qEh2gXQEOsFS3r5AWRKwKBgQC8
a8QBxaP36gQjNc9Zmwqm9zmOysQuu8nQQM8GOr9RxSFl2X0PEVx6RUyokgo3xgHQ
38qgWAxbXgeWuNcaBWuH+OgvgFYUViz5U8GqjfDvs2lzTNttIMw63ylNUw2SiMPg
OzvunwuEu/80Wy0Kcy37zcMhoSgF1a6vjC+EV3uNGwKBgFdMgUcI9qlps+yjK2vJ
zpg8YNiScpzDQC78fgw1VhPtYFDlxfLJ+gGVWc5IE3PFe/egCKSd0yVbun3LVE1/
/DMkCySIcfo1hKHcLKnpg1M+gqLb5G02aCmjWyCC8sV9EYYT6dsuBR9xX4lRXL5F
nkuGiPQtUMzAXjyKlQnwtf9H
-----END PRIVATE KEY-----`

	CreatorPassword  = ""
	ChainTrackerId   = "22AErGxiofaGmPrMYPzq54zvBuGW6RAWb8e8MwWD7BAC"
	ChainTrackerAddr = "tracker.tethyscore.io"
	ChainTrackerPort = "80"

	AnonUserPKPem = `MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEzAChkyGokyyQRCq9Gls+ynmolSS4
NIBRCRBlNi+GIZqmpGySeroopjA+BUbIk6HuLsEllICMMqFfQirl3/Tg1w==`
	AnonUserPKDer = `MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEzAChkyGokyyQRCq9Gls+ynmolSS4NIBRCRBlNi+GIZqmpGySeroopjA+BUbIk6HuLsEllICMMqFfQirl3/Tg1w==`
	AnonUserSK    = `-----BEGIN PRIVATE KEY-----
MIGTAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBHkwdwIBAQQgXMUB59Nii+UdOiX3
Rkr7ean14G05Iipd/GHAxuGIiamgCgYIKoZIzj0DAQehRANCAATMAKGTIaiTLJBE
Kr0aWz7KeaiVJLg0gFEJEGU2L4YhmqakbJJ6uiimMD4FRsiToe4uwSWUgIwyoV9C
KuXf9ODX	
-----END PRIVATE KEY-----`

	WorldSign = `v+xeZN1nlVitwGRXhePqVUAjCpqb7q8CQtX3cZJUvUg5wPnBG4GFzbGPG9WZNj+k5iBIMYlOKETzP59eBmGy5tAjWTWle/tOMuCEBaIX//TMYv5WBaxTG4ZIyVX2lIF+M0YphCzxyoBJoQIqgnrlzzwfY6mSlL/Er90ntmO1Jc6jsaKZGmbaiU+lv/425BscGXZDqKSulyGSAqImo26thgDNWL5qeQn8svHvlF79E9IANYxqLaQRP/qsoFVefhtzhmWsHYavp4bmTxjJP8cx91z+3SWqyalCzBx8Ff2+ap+aiMEtoyc7YF3w+kQLKBb1odcIqZDv2sQeBEXBDVrdeg==`
)

// This test codes are not for testing itself but generating sample or actual chian id and signatures for world and chain
// You can replace const values above and make the values you're looking for
func TestWorldChain(t *testing.T) {
	cert, err := GetCertificate(CreatorCert)
	assert.Nil(t, err, "cannot parse certificate of creator")

	creatorIdRaw := sha256.Sum256([]byte(cert.Subject.CommonName))
	creatorIdB58 := base58.Encode(creatorIdRaw[:])
	fmt.Println("creator id: ", creatorIdB58)

	bBuilder := bytesbuilder.NewBuilder()
	bBuilder.Append(WorldName)
	bBuilder.Append(ChainName)
	cTime, err := strconv.ParseInt(ChainTime, 10, 64)
	assert.Nil(t, err, "cannot parse timestamp of chain")
	bBuilder.Append(cTime)
	bBuilder.Append(ChainCustomContract)
	bBuilder.Append(ChainOracle)
	bBuilder.Append(ChainDym)
	bBuilder.Append(ChainScript)
	bBuilder.AppendBase58(creatorIdB58)

	chainIdRaw := sha256.Sum256(bBuilder.GetBytes())
	chainIdB64 := base64.StdEncoding.EncodeToString(chainIdRaw[:])
	// chainIdB58 := base58.Encode(chainIdRaw[:])
	fmt.Println("chain id: ", chainIdB64)

	bBuilder.Clear()

	bBuilder.AppendBase64(chainIdB64)
	bBuilder.AppendBase58(ChainTrackerId)
	bBuilder.Append(ChainTrackerAddr)
	trackerPort, err := strconv.Atoi(ChainTrackerPort)
	assert.Nil(t, err, "cannot parse port of tracker")
	bBuilder.Append(int32(trackerPort))
	chainMsgRaw := bBuilder.GetBytes()

	skey, err := GetPrivateKey(CreatorKey, "")
	assert.Nil(t, err, "cannot parse secret key of creator")

	chainSig, err := DoSign(chainMsgRaw, skey)
	assert.Nil(t, err, "cannot make signature of chain creator")
	fmt.Println("chain sig: ", base64.StdEncoding.EncodeToString(chainSig))
	chainSignRes := DoVerify(chainMsgRaw, chainSig, cert.PublicKey)
	assert.True(t, chainSignRes, "cannot verify signature of chain")

	bBuilder.Clear()

	bBuilder.Append(WorldName)
	bBuilder.Append(cTime)
	bBuilder.Append(WorldKeyCcyName)
	keyCcyAmount, err := strconv.ParseInt(WorldKeyCcyAmount, 10, 64)
	assert.Nil(t, err, "cannot parse key currency amount of chain")
	bBuilder.Append(keyCcyAmount)
	bBuilder.Append(WorldAnonymousUser)
	bBuilder.AppendBase64(chainIdB64)
	bBuilder.AppendBase58(creatorIdB58)
	bBuilder.Append(CreatorCert)
	bBuilder.Append(WorldAuthorityId)
	msgRaw := bBuilder.GetBytes()

	signRaw, err := DoSign(msgRaw, skey)
	assert.Nil(t, err, "cannot make signature of world creator")

	fmt.Println("world sig: ", base64.StdEncoding.EncodeToString(signRaw))
	worldSignRes := Verify(msgRaw, signRaw, CreatorCert)
	assert.True(t, worldSignRes, "cannot verify signature of world")
}

func TestAnonymousUser(t *testing.T) {
	testMsg := []byte("hello")
	skey, err := GetPrivateKey(AnonUserSK, "")
	fmt.Println(reflect.TypeOf(skey))
	assert.Nil(t, err, "What happened")
	signature, err := Sign(testMsg, AnonUserSK, "")
	assert.Nil(t, err, "1111")
	// pBlock, _ := pem.Decode([]byte(AnonUserPKPem))
	// pBlockPemRaw := pBlock.Bytes
	// // derTrim := strings.ReplaceAll(AnonUserPKDer, "\n", "")
	// // derTrim = strings.ReplaceAll(derTrim, "\r", "")

	// derRaw, err := base64.StdEncoding.DecodeString(AnonUserPKDer)

	// fmt.Println(pBlockPemRaw)
	// fmt.Println(derRaw)

	// pk, err := GetPublicKey(derRaw)
	// fmt.Println(reflect.TypeOf(pk))
	success := Verify(testMsg, signature, AnonUserPKDer)
	assert.True(t, success, "2222")

	hashed := sha256.Sum256([]byte(AnonUserPKPem))
	fmt.Println("user id: ", base58.Encode(hashed[:]))

}

func TestWorldSign(t *testing.T) {
	bBuilder := bytesbuilder.NewBuilder()
	bBuilder.Append(WorldName)
	wTime, err := strconv.ParseInt(WorldTime, 10, 64)
	assert.Nil(t, err, "cannot parse timestamp of world")
	bBuilder.Append(wTime)
	bBuilder.Append(WorldKeyCcyName)
	keyCcyAmount, err := strconv.ParseInt(WorldKeyCcyAmount, 10, 64)
	assert.Nil(t, err, "cannot parse key currency amount of world")
	bBuilder.Append(keyCcyAmount)
	bBuilder.Append(WorldAnonymousUser)
	bBuilder.AppendBase64(WorldEdenChainId)
	bBuilder.AppendBase58(WorldCreatorId)
	bBuilder.Append(CreatorCert)
	bBuilder.Append(WorldAuthorityId)
	msgRaw := bBuilder.GetBytes()

	signRaw, err := base64.StdEncoding.DecodeString(WorldSign)
	assert.Nil(t, err, "cannot decode signature of world")

	success := Verify(msgRaw, signRaw, CreatorCert)
	assert.True(t, success, "should be true")
}
