package util

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
)

func Sha1(v, l string) string {
	key := []byte(l)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(v))
	return hex.EncodeToString(mac.Sum(nil))
}
