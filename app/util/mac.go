package util

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
)

func CreateSign(message, key []byte) string {
	hmacSha1 := hmac.New(sha1.New, key)
	mac := hmacSha1.Sum(message)
	sign := hex.EncodeToString(mac)
	return base64.RawURLEncoding.EncodeToString([]byte(sign))
}
