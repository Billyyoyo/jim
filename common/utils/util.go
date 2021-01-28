package utils

import (
	"bytes"
	"crypto"
	"crypto/cipher"
	"crypto/des"
	"crypto/hmac"
	"crypto/md5"
	randCrypto "crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	randNum "math/rand"
	"time"
)

const (
	DEFAULT_LETTERS string = "AB0CD1EFX2GH3IJ4WKL5MN6YOP7QZR8ST9UV"
	DES_KEY         string = "jim_hj80"
	HMAC_KEY        string = "jim_hj80"
)

var PRIVATE_KEY = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQClMDAT31yB65zoccklBR5nOrXemBOs7mke1o5+5NGsW7vbWmVa
Dz5YzSmRCNkC9BedOah6AS5Xi4KX9oznbHiHupmouFKNzI5QWrh9z0vyr5ef+zn3
G6K0jyVP11TZyR2YSKvkbH0mzdu2d49X4hH5ZNc00KLKndNqzOyQQnZqJwIDAQAB
AoGAGRBuuxUxHCV78WkSdCOKsW8fGV9J3Ptvx9YWWPqvNc/VeTsGFdwqQZ8fp3oI
M4KF8r1E6v4y4eWxMw2d4595yh/lLvpbJogoq3zXIthB55P9oGHGV+jOQPXz64Ya
uVpFQ3jIdkO06L6E7/fesvcIrc1gdSjO/c4lgVu/6c2KQwECQQDbVgbvr8f6n8Bh
jkssTLFhai+YoRM3C3RJhm5/31ROhhKTtX1N3BMygXBvGXBeWgEpXQ8ECGwXO+je
ZN/lgNhXAkEAwM0IJ4aIxc0apNzEd5/nJumRWmRUKM8E2V6VtGOeEI/BEl0ysz7t
n4gIvbvbAn345naUxoYD+t7Ht7a1k9gasQJAYlnZ9GJjDsvRjS0sIiolo+PkgdFA
d39IXqvMIsS23hsae2d3T5FufkgybW7/xx8exDh5QjqwlV6E1ixvhU7YMwJAQpqR
5pWjSjHAspNRi8HBqL+nZwKh0Dc0BaOXM+n2AOKoYB+yFBn5HNNxsZnj3siF45ez
baF+Xnv3oo+LyrctAQJBAIL4Ww3b4054V4iCNMb7poNuodFmcxIpgqBKy0NsKk2B
Svk0LnCxPIMCr3LeoFXAkre+SzOUnyQuGjcgMwJUcQc=
-----END RSA PRIVATE KEY-----`)

var PUBLIC_KEY = []byte(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQClMDAT31yB65zoccklBR5nOrXe
mBOs7mke1o5+5NGsW7vbWmVaDz5YzSmRCNkC9BedOah6AS5Xi4KX9oznbHiHupmo
uFKNzI5QWrh9z0vyr5ef+zn3G6K0jyVP11TZyR2YSKvkbH0mzdu2d49X4hH5ZNc0
0KLKndNqzOyQQnZqJwIDAQAB
-----END PUBLIC KEY-----`)

func GetCurrentMS() int64 {
	return time.Now().UnixNano() / 1e6
}

func IF(cond bool, val1, val2 interface{}) interface{} {
	if cond {
		return val1
	}
	return val2
}

func RandomStr(n int, allowedChars ...string) string {
	var letters string

	if len(allowedChars) == 0 {
		letters = DEFAULT_LETTERS
	} else {
		letters = allowedChars[0]
	}
	b := make([]byte, n)
	for i := range b {
		randNum.Seed(time.Now().UnixNano())
		b[i] = letters[randNum.Intn(len(letters))]
	}

	return string(b)
}

func SignRsa(privateKey []byte, sourceData []byte) (string, error) {
	block, _ := pem.Decode(privateKey)
	priKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	myHash := sha256.New()
	myHash.Write(sourceData)
	hashRes := myHash.Sum(nil)
	res, err := rsa.SignPKCS1v15(randCrypto.Reader, priKey, crypto.SHA256, hashRes)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(res), nil
}

func VerifyRsa(publicKey []byte, sourceData []byte, signedData string) error {
	ciphertext, _ := base64.URLEncoding.DecodeString(signedData)
	block, _ := pem.Decode(publicKey)
	publicInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	pubKey := publicInterface.(*rsa.PublicKey)
	mySha := sha256.New()
	mySha.Write(sourceData)
	res := mySha.Sum(nil)
	err = rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, res, ciphertext)
	if err != nil {
		return err
	}
	return nil
}
func EncryptRSA(origData []byte, publicKey []byte) ([]byte, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(randCrypto.Reader, pub, origData)
}

func DecryptRsa(ciphertext []byte, privateKey []byte) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(randCrypto.Reader, priv, ciphertext)
}

func padding(src []byte, blocksize int) []byte {
	n := len(src)
	padnum := blocksize - n%blocksize
	pad := bytes.Repeat([]byte{byte(padnum)}, padnum)
	dst := append(src, pad...)
	return dst
}

func unpadding(src []byte) []byte {
	n := len(src)
	unpadnum := int(src[n-1])
	dst := src[:n-unpadnum]
	return dst
}

func EncryptDES(src []byte, key []byte) []byte {
	block, _ := des.NewCipher(key)
	src = padding(src, block.BlockSize())
	blockmode := cipher.NewCBCEncrypter(block, key)
	blockmode.CryptBlocks(src, src)
	return []byte(base64.StdEncoding.EncodeToString(src))
}

func DecryptDES(b64 []byte, key []byte) []byte {
	src, _ := base64.StdEncoding.DecodeString(string(b64))
	block, _ := des.NewCipher(key)
	blockmode := cipher.NewCBCDecrypter(block, key)
	blockmode.CryptBlocks(src, src)
	src = unpadding(src)
	return src
}

func SignHmac(origin []byte, key []byte) []byte {
	tool := hmac.New(sha1.New, key)
	tool.Write(origin)
	return tool.Sum(nil)
}

func Md5sum(text string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(text)))
}
