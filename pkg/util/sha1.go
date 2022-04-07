package util

import (
	"crypto/hmac"
	"crypto/sha1"
)

func SHA1(v, l string) []byte {
	key := []byte(l)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(v))
	return mac.Sum(nil)
}

