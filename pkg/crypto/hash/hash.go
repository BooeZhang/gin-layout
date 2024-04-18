package hash

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

// MD5 md5 hash
func MD5(b []byte) string {
	h := md5.New()
	_, _ = h.Write(b)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// MD5String md5 hash
func MD5String(s string) string {
	return MD5([]byte(s))
}

// SHA1 sha1 hash
func SHA1(b []byte) string {
	h := sha1.New()
	_, _ = h.Write(b)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// SHA1String sha1 hash
func SHA1String(s string) string {
	return SHA1([]byte(s))
}

// ValidHmac 校验 hmac
func ValidHmac(data []byte, key []byte, hmacClient string) (ok bool) {
	return HmacStr(data, key) == hmacClient
}

func HmacStr(data []byte, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write(data)

	_hmac := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return _hmac
}
