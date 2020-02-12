package moc

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/btcsuite/btcutil/base58"
	"github.com/lokks307/go-util/bytesbuilder"
	"github.com/stretchr/testify/assert"
)

const (
	WorldName          = "TESTWORLD"
	WorldTime          = "1581500000"
	WorldAnonymousUser = true
	WorldKeyCcyName    = "TYS"
	WorldKeyCcyAmount  = "100000000"
	WorldAuthorityId   = "*"

	ChainName           = "TESTCHAIN"
	ChainTime           = "1581500000"
	ChainCustomContract = false
	ChainOracle         = false
	ChainDym            = false
	ChainScript         = false

	CreatorCN   = "H2KpjFrdhat2uHy277EKsyfKV3fz1PxS7XqgQQNGGL1d"
	CreatorCert = `-----BEGIN CERTIFICATE-----
MIICKzCCAdECFHQiXWrVH4tU2qe3DVa994v9pC3mMAoGCCqGSM49BAMCMIGXMQsw
CQYDVQQGEwJLUjEQMA4GA1UECAwHSW5jaGVvbjETMBEGA1UEBwwKTmFtZG9uZy1n
dTEVMBMGA1UECgwMTG9ra3MzMDcgSW5jMRMwEQYDVQQLDApUZXN0IHVzZXIxMTUw
MwYDVQQDDCxIMktwakZyZGhhdDJ1SHkyNzdFS3N5ZktWM2Z6MVB4UzdYcWdRUU5H
R0wxZDAeFw0yMDAyMTEwODQyMTdaFw0zMDAyMDgwODQyMTdaMIGXMQswCQYDVQQG
EwJLUjEQMA4GA1UECAwHSW5jaGVvbjETMBEGA1UEBwwKTmFtZG9uZy1ndTEVMBMG
A1UECgwMTG9ra3MzMDcgSW5jMRMwEQYDVQQLDApUZXN0IHVzZXIxMTUwMwYDVQQD
DCxIMktwakZyZGhhdDJ1SHkyNzdFS3N5ZktWM2Z6MVB4UzdYcWdRUU5HR0wxZDBZ
MBMGByqGSM49AgEGCCqGSM49AwEHA0IABJM1pWbB35WwODIqWpOg9MxsLqsK/9Et
hm7woZbF2gZIM0Xwcxc5QFhzpXDPXFQlbaFCC2gKDrrtBoLDTetm3QQwCgYIKoZI
zj0EAwIDSAAwRQIgQ5WXgEz9/t5k7Cb/fmSjzr2HoBioIs62PoZ/9QLDdpcCIQDP
tE8MIQ1Hc4j7pw2I3T0vnm9BEBOTG83s0HAynp0a7g==
-----END CERTIFICATE-----`

	CreatorKey = `-----BEGIN PRIVATE KEY-----
MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQgwgcBEKLAdJ+MN/2j
/7XCLdaSBqGb1UHHZJIbXZOzrS+hRANCAASTNaVmwd+VsDgyKlqToPTMbC6rCv/R
LYZu8KGWxdoGSDNF8HMXOUBYc6Vwz1xUJW2hQgtoCg667QaCw03rZt0E
-----END PRIVATE KEY-----`

	CreatorPassword  = ""
	ChainTrackerId   = "41fCjrkvCHjEW7VbVUBH2oww1TkA3ET8KrZSVtDYzER6"
	ChainTrackerAddr = "tracker.example.io"
	ChainTrackerPort = "80"

	AnonUserPKPem = `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEzAChkyGokyyQRCq9Gls+ynmolSS4
NIBRCRBlNi+GIZqmpGySeroopjA+BUbIk6HuLsEllICMMqFfQirl3/Tg1w==
-----END PUBLIC KEY-----`
	AnonUserPKDer = `MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEzAChkyGokyyQRCq9Gls+ynmolSS4NIBRCRBlNi+GIZqmpGySeroopjA+BUbIk6HuLsEllICMMqFfQirl3/Tg1w==`
	AnonUserSK    = `-----BEGIN PRIVATE KEY-----
MIGTAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBHkwdwIBAQQgXMUB59Nii+UdOiX3
Rkr7ean14G05Iipd/GHAxuGIiamgCgYIKoZIzj0DAQehRANCAATMAKGTIaiTLJBE
Kr0aWz7KeaiVJLg0gFEJEGU2L4YhmqakbJJ6uiimMD4FRsiToe4uwSWUgIwyoV9C
KuXf9ODX	
-----END PRIVATE KEY-----`
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
	chainIdB58 := base58.Encode(chainIdRaw[:])
	fmt.Println("chain id: ", chainIdB58)

	bBuilder.Clear()

	bBuilder.AppendBase58(chainIdB58)
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
	bBuilder.AppendBase58(chainIdB58)
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
	pBlock, _ := pem.Decode([]byte(AnonUserPKPem))
	pBlockPemRaw := pBlock.Bytes
	// derTrim := strings.ReplaceAll(AnonUserPKDer, "\n", "")
	// derTrim = strings.ReplaceAll(derTrim, "\r", "")
	derRaw, err := base64.StdEncoding.DecodeString(AnonUserPKDer)

	fmt.Println(pBlockPemRaw)
	fmt.Println(derRaw)

	pk, err := GetPublicKey(derRaw)
	fmt.Println(reflect.TypeOf(pk))
	success := DoVerify(testMsg, signature, pk)
	assert.True(t, success, "2222")

	hashed := sha256.Sum256([]byte(AnonUserPKPem))
	fmt.Println("user id: ", base58.Encode(hashed[:]))

}
