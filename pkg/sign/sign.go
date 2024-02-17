package sign

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

// ValidHmacSign 校验 hmac
func ValidHmacSign(data []byte, key []byte, hmacClient string) (ok bool) {
	return HmacSign(data, key) == hmacClient
}

func HmacSign(data []byte, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write(data)

	_hmac := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return _hmac
}

func MD5String(str string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}
