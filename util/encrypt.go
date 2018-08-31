package util

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
)

// HmacSha1Base64 执行HmacSha1Base64加密
func HmacSha1Base64(data, key string) string {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
