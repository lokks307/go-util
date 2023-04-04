package moc

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"math/big"
	"strings"

	"github.com/lokks307/pkcs8"
	"software.sslmate.com/src/go-pkcs12"
)

// MOdern Cryptography library for TethysCore

const BeginPEMFormat = "-----BEGIN"
const EndPEMFormat = "-----"
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
	dataB64 = strings.TrimSpace(dataB64)

	if IsPEMFormat([]byte(dataB64)) {
		pemBlk, _ := pem.Decode([]byte(dataB64))
		if pemBlk == nil {
			return nil, errors.New("Failed to decode PEM format")
		}

		var raw []byte

		if x509.IsEncryptedPEMBlock(pemBlk) {
			decryptRaw, decryptErr := x509.DecryptPEMBlock(pemBlk, []byte(password))
			if decryptErr != nil {
				return nil, errors.New("Failed to decrypt PEM, " + decryptErr.Error())
			} else {
				raw = decryptRaw
			}
		} else {
			raw = pemBlk.Bytes
		}

		switch pemBlk.Type {
		case "RSA PRIVATE KEY":
			key, err = x509.ParsePKCS1PrivateKey(raw)
			return
		case "EC PRIVATE KEY":
			key, err = x509.ParseECPrivateKey(raw)
			return
		case "ENCRYPTED PRIVATE KEY":
			key, err = pkcs8.ParsePKCS8PrivateKey(raw, []byte(password))
			return
		default:
			key, err = x509.ParsePKCS8PrivateKey(raw)
			return
		}
	} else {
		// It is not PEM format. Probably base encoded DER format
		// WARN: encrypted DER is not supported
		dataB64 = strings.ReplaceAll(dataB64, "\n", "")
		dataB64 = strings.ReplaceAll(dataB64, "\r", "")

		dataB64Raw, decodeErr := base64.StdEncoding.DecodeString(dataB64)
		if err != nil {
			return nil, errors.New("Failed to decode DER base64, " + decodeErr.Error())
		}

		// cannot check pkcs1 or pkcs8 type in der format. so try pkcs8 first and try pkcs1 again
		key, err = x509.ParsePKCS8PrivateKey(dataB64Raw)
		if err == nil {
			return
		}

		key, err = x509.ParseECPrivateKey(dataB64Raw)
		if err == nil {
			return
		}

		key, err = x509.ParsePKCS1PrivateKey(dataB64Raw)
		if err == nil {
			return
		}

		return nil, errors.New("Failed to parse private key")
	}
}

func IsPEMFormat(in []byte) bool {
	return bytes.HasPrefix(in, []byte(BeginPEMFormat)) && bytes.HasSuffix(in, []byte(EndPEMFormat))
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
