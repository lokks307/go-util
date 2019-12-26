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

func GetCertificate(dataB64 string) (*x509.Certificate, error) {
	data := ParseDataToDer(dataB64)

	if data == nil {
		return nil, errors.New("Can't decode data")
	}

	var cert *x509.Certificate
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
	cert, parseErr := GetCertificate(dataB64)

	if parseErr != nil {
		dataRaw := ParseDataToDer(dataB64)
		pubKey, err := x509.ParsePKIXPublicKey(dataRaw)

		if err != nil {
			return nil, err
		}
		return pubKey, nil
	}

	return cert.PublicKey, nil
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
		var rBigIntRaw []byte
		var sBigIntRaw []byte
		if len(r.Bytes()) < 32 { // FIXME: what if we have to support other eliptical curve like P-192, P-521?
			for i := 0; i < 32-len(r.Bytes()); i++ {
				rBigIntRaw[i] = 0x00
			}
		}
		rBigIntRaw = append(rBigIntRaw, r.Bytes()...)
		if len(s.Bytes()) < 32 {
			for i := 0; i < 32-len(s.Bytes()); i++ {
				sBigIntRaw[i] = 0x00
			}
		}
		sBigIntRaw = append(sBigIntRaw, s.Bytes()...)
		signature = append(rBigIntRaw, sBigIntRaw...)
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

func DoVerify(msg, sigBytes []byte, publicKey interface{}) bool {
	var result bool

	hashed := sha256.Sum256(msg)

	switch pubKey := publicKey.(type) {
	case *rsa.PublicKey:
		err := VerifyRSASign(hashed[:], sigBytes, pubKey)
		if err != nil {
			result = false
		} else {
			result = true
		}
	case *ecdsa.PublicKey:
		if len(sigBytes) > 64 { // signature in DER format
			var ecdsaInts []*big.Int

			_, err := asn1.Unmarshal(sigBytes, &ecdsaInts)

			if err != nil {
				return false
			}

			result = ecdsa.Verify(pubKey, hashed[:], ecdsaInts[0], ecdsaInts[1])
		} else {
			halfSigLen := len(sigBytes) / 2
			r := new(big.Int)
			r.SetBytes(sigBytes[:halfSigLen])

			s := new(big.Int)
			s.SetBytes(sigBytes[halfSigLen:])

			result = ecdsa.Verify(pubKey, hashed[:], r, s)
		}
	case ed25519.PublicKey:
		result = ed25519.Verify(pubKey, msg, sigBytes)
	default:
		result = false
	}

	return result
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

func VerifyRSASign(hashMsgRaw, signRaw []byte, rsaPubKey *rsa.PublicKey, mode string) error {
	var err error
	err = nil

	if len(mode) == 0 || mode == "pkcs1v15" { // default
		err = rsa.VerifyPKCS1v15(rsaPubKey, crypto.SHA256, hashMsgRaw, signRaw)
	} else if mode == "pss" {
		err = rsa.VerifyPSS(rsaPubKey, crypto.SHA256, hashMsgRaw, signRaw, nil)
	}

	return err
}
