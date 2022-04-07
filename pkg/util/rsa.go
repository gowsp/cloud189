package util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
)

func Key(key string) []byte {
	return []byte("-----BEGIN PUBLIC KEY-----\n" + key + "\n-----END PUBLIC KEY-----")
}

func RsaEncrypt(key []byte, data []byte) ([]byte, error) {
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, data)
}

func Encrypt(req *http.Request, secret string) {
	// sha1(SessionKey=相应的值&Operate=相应值&RequestURI=相应值&Date=相应的值”, SessionSecret)
	params := req.Header
	d := fmt.Sprintf("SessionKey=%s&Operate=%s&RequestURI=%s&Date=%s", params.Get("SessionKey"),
		req.Method, req.URL.Path, params.Get("Date"))
	val := SHA1(d, secret)
	req.Header.Set("Signature", hex.EncodeToString(val))
	req.Header.Set("X-Request-ID", Random("xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx"))
}
