package util

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
)

const (
	PKCS1 int64 = iota
	PKCS8

	// 私钥 PEMBEGIN 开头
	PEMBEGIN = "-----BEGIN RSA PRIVATE KEY-----\n"
	// 私钥 PEMEND 结尾
	PEMEND = "\n-----END RSA PRIVATE KEY-----"
	// 公钥 PEMBEGIN 开头
	PUBPEMBEGIN = "-----BEGIN PUBLIC KEY-----\n"
	// 公钥 PEMEND 结尾
	PUBPEMEND = "\n-----END PUBLIC KEY-----"
)

func Sha256WithRsa(privateRaw string, msg string) (string, error) {
	privateRaw = strings.Trim(privateRaw, "\n")
	if !strings.HasPrefix(privateRaw, "-----BEGIN RSA PRIVATE KEY-----") {
		privateRaw = fmt.Sprintf("%s\n%s\n%s", "-----BEGIN RSA PRIVATE KEY-----", privateRaw, "-----END RSA PRIVATE KEY-----")
	}

	blockPri, _ := pem.Decode([]byte(privateRaw))
	if blockPri == nil {
		return "", fmt.Errorf("blockPri is nil")
	}

	rsaPri, e := genPriKey(blockPri.Bytes, PKCS8)
	if e != nil {
		panic(e)
	}

	h := sha256.New()
	h.Write([]byte(msg))
	d := h.Sum(nil)

	signature, err := rsa.SignPKCS1v15(rand.Reader, rsaPri, crypto.SHA256, d)
	if err != nil {
		return "", errors.New("错误")
	}
	encodedSig := base64.StdEncoding.EncodeToString(signature)
	return encodedSig, nil
}

func Rsa2PubSign(publicKey, signContent, sign string) bool {
	hashed := sha256.Sum256([]byte(signContent))
	pubKey, err := ParsePublicKey(publicKey)
	if err != nil {
		log.Error(err, "rsa2 public check sign failed.")
		return false
	}
	sig, _ := base64.StdEncoding.DecodeString(sign)
	err = rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hashed[:], sig)
	if err != nil {
		log.Error(err, "rsa2 public check sign failed.")
		return false
	}
	return true
}

func genPriKey(privateKey []byte, privateKeyType int64) (*rsa.PrivateKey, error) {
	var priKey *rsa.PrivateKey
	var err error
	switch privateKeyType {
	case PKCS1:
		{
			priKey, err = x509.ParsePKCS1PrivateKey([]byte(privateKey))
			if err != nil {
				return nil, err
			}
		}
	case PKCS8:
		{
			prkI, err := x509.ParsePKCS8PrivateKey([]byte(privateKey))
			if err != nil {
				return nil, err
			}
			priKey = prkI.(*rsa.PrivateKey)
		}
	default:
		{
			return nil, fmt.Errorf("unsupport private key type")
		}
	}
	return priKey, nil
}

func ParsePublicKey(publicKey string) (*rsa.PublicKey, error) {
	publicKey = FormatPublicKey(publicKey)
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return nil, errors.New("公钥信息错误！")
	}
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pubKey.(*rsa.PublicKey), nil
}

// FormatPublicKey 组装公钥
func FormatPublicKey(publicKey string) string {
	if !strings.HasPrefix(publicKey, PUBPEMBEGIN) {
		publicKey = PUBPEMBEGIN + publicKey
	}
	if !strings.HasSuffix(publicKey, PUBPEMEND) {
		publicKey = publicKey + PUBPEMEND
	}
	return publicKey
}
