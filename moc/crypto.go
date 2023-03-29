package moc

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/hex"
	"errors"
	"math/big"

	"github.com/lokks307/pkcs8"
	"software.sslmate.com/src/go-pkcs12"
)

// MOdern Cryptography library for TethysCore

const BeginPhraseCertificate = "-----BEGIN CERTIFICATE-----"

func GetCertificateOrPublicKey(dataB64 string) (interface{}, error) {
	if dataB64 == "" {
		return nil, errors.New("Empty input")
	}

	if BeginPhraseCertificate == dataB64[0:27] { // this is PEM certificate
		if pBlock, err := DecodePEM(dataB64); err == nil {
			if cert, err := x509.ParseCertificate(pBlock.Bytes); err == nil {
				return cert, nil
			} else {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	dataRaw := ParseDataToDer(dataB64)
	if dataRaw == nil { // wrong encoded data
		return nil, errors.New("Cannot decode Base64 format")
	}

	if cert, err := x509.ParseCertificate(dataRaw); err == nil {
		return cert, nil
	}

	if pubKey, err := x509.ParsePKIXPublicKey(dataRaw); err == nil {
		return pubKey, nil
	}

	return nil, errors.New("Cannot parse data to x509 or public key")
}

func GetCertificate(dataB64 string) (*x509.Certificate, error) {
	data := ParseDataToDer(dataB64)

	if data == nil {
		return nil, errors.New("Can't decode data")
	}

	cert, parseErr := x509.ParseCertificate(data)

	if parseErr != nil {
		return nil, parseErr
	}

	return cert, nil
}

func DecodePFX(der []byte, password string) (privateKey interface{}, certificate *x509.Certificate, err error) {
	privKey, cert, err := pkcs12.Decode(der, password)
	if err != nil {
		return nil, nil, err
	}

	return privKey, cert, nil
}

func DecodePFXB64(pfxDataB64 string, password string) (privateKey interface{}, certificate *x509.Certificate, err error) {
	res := ParseDataToDer(pfxDataB64)

	return DecodePFX(res, password)
}

func GetPrivateKey(dataB64, password string) (key interface{}, err error) {
	derBytes := ParseDataToDer(dataB64, password)

	// cannot check pkcs1 or pkcs8 type in der format. so try pkcs8 first and try pkcs1 again
	privKey, _, err := pkcs8.ParsePrivateKey(derBytes, []byte(password))
	if err != nil {
		pkcs1Key, parseErr := x509.ParsePKCS1PrivateKey(derBytes)

		if parseErr != nil {
			return nil, parseErr
		}

		return pkcs1Key, nil
	}

	return privKey, nil
}

func GetPublicKey(dataB64 string) (public interface{}, err error) {
	outItf, err := GetCertificateOrPublicKey(dataB64)
	if err != nil {
		return nil, err
	}

	switch output := outItf.(type) {
	case *x509.Certificate:
		return output.PublicKey, nil
	default:
		return output, nil
	}
}

func makeSigFromBigInt(r, s *big.Int, lenField int) []byte {
	var sig []byte
	var rBigIntRaw []byte
	if len(r.Bytes()) < lenField {
		for i := 0; i < lenField-len(r.Bytes()); i++ {
			rBigIntRaw = append(rBigIntRaw, 0x00)
		}
	}
	rBigIntRaw = append(rBigIntRaw, r.Bytes()...)

	var sBigIntRaw []byte
	if len(s.Bytes()) < lenField {
		for i := 0; i < lenField-len(s.Bytes()); i++ {
			sBigIntRaw = append(sBigIntRaw, 0x00)
		}
	}
	sBigIntRaw = append(sBigIntRaw, s.Bytes()...)

	sig = append(rBigIntRaw, sBigIntRaw...)

	return sig
}

func DoSign(msg []byte, key interface{}) ([]byte, error) {
	rng := rand.Reader
	var signature []byte
	var err error
	hashed := sha256.Sum256(msg)
	err = nil

	switch privKey := key.(type) {
	case *rsa.PrivateKey:
		signature, err = rsa.SignPSS(rng, privKey, crypto.SHA256, hashed[:], nil)
	case *ecdsa.PrivateKey:
		var r *big.Int
		var s *big.Int
		r, s, err = ecdsa.Sign(rng, privKey, hashed[:])
		if err == nil {
			bitLen := privKey.Curve.Params().BitSize
			signature = makeSigFromBigInt(r, s, int(bitLen/8))
		}
	case ed25519.PrivateKey:
		signature = ed25519.Sign(privKey, msg)
	default:
		signature = nil
		err = errors.New("Unsupported type of crypto method")
	}

	return signature, err
}

func Sign(msg []byte, skeyPem, pass string) ([]byte, error) {
	privateKey, err := GetPrivateKey(skeyPem, pass)
	if err != nil {
		return nil, err
	}

	return DoSign(msg, privateKey)
}

func DoVerify(msgRaw, sigRaw []byte, publicKey interface{}) bool {

	switch pubKey := publicKey.(type) {
	case *rsa.PublicKey:
		return VerifySignatureRSAPSS(msgRaw, sigRaw, pubKey)
	case *ecdsa.PublicKey:
		return VerifySignatureECDSA(msgRaw, sigRaw, pubKey)
	case ed25519.PublicKey:
		return VerifySignatureEDDSA(msgRaw, sigRaw, pubKey)
	default:
		return false
	}
}

func Verify(msg, sigBytes []byte, certPem string) bool {
	publicKey, err := GetPublicKey(certPem)

	if err != nil {
		return false
	}

	return DoVerify(msg, sigBytes, publicKey)
}

func VerifyFromHexPubKey(msg, sigBytes []byte, hexStr string) bool {
	derRaw, err := hex.DecodeString(hexStr)
	if err != nil {
		return false
	}

	cert, certErr := x509.ParseCertificate(derRaw)
	if certErr != nil {
		pubKey, pubErr := x509.ParsePKIXPublicKey(derRaw)
		if pubErr != nil {
			return false
		}
		return DoVerify(msg, sigBytes, pubKey)
	}
	return DoVerify(msg, sigBytes, cert.PublicKey)
}

func VerifySignatureEDDSA(msgRaw, sigRaw []byte, pubKey interface{}) bool {
	edKey, ok := pubKey.(ed25519.PublicKey)
	if !ok {
		return false
	}

	return ed25519.Verify(edKey, msgRaw, sigRaw)
}

func VerifySignatureECDSA(msgRaw, sigRaw []byte, pubKey interface{}) bool {
	ecKey, ok := pubKey.(*ecdsa.PublicKey)
	if !ok {
		return false
	}

	hashed := sha256.Sum256(msgRaw)

	var esig struct {
		R, S *big.Int
	}

	_, err := asn1.Unmarshal(sigRaw, &esig)
	if err != nil { // signature is not der format
		halfSigLen := len(sigRaw) / 2
		r := new(big.Int)
		r.SetBytes(sigRaw[:halfSigLen])

		s := new(big.Int)
		s.SetBytes(sigRaw[halfSigLen:])

		return ecdsa.Verify(ecKey, hashed[:], r, s)
	} else {
		return ecdsa.Verify(ecKey, hashed[:], esig.R, esig.S)
	}
}

func VerifySignatureRSAPSS(msgRaw, sigRaw []byte, pubKey interface{}) bool {
	rsaKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return false
	}

	// pssOpts := rsa.PSSOptions{SaltLength: 20, Hash: 0}
	hashed := sha256.Sum256(msgRaw)
	err := rsa.VerifyPSS(rsaKey, crypto.SHA256, hashed[:], sigRaw, nil)
	if err != nil {
		return false
	} else {
		return true
	}
}

func VerifySignatureRSAPKCS1(msgRaw, sigRaw []byte, pubKey interface{}) bool {
	rsaKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return false
	}

	hashed := sha256.Sum256(msgRaw)
	err := rsa.VerifyPKCS1v15(rsaKey, crypto.SHA256, hashed[:], sigRaw)
	if err != nil {
		return false
	} else {
		return true
	}
}
