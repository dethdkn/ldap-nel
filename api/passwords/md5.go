package passwords

import (
	"crypto/md5"
	"crypto/subtle"
	"encoding/base64"
	"strings"
)

func MD5(password string) string {
	sum := md5.Sum([]byte(password))
	return "{MD5}" + base64.StdEncoding.EncodeToString(sum[:])
}

func VerifyMD5(password, hashed string) bool {
	if !strings.HasPrefix(hashed, "{MD5}") {
		return false
	}
	b64 := strings.TrimPrefix(hashed, "{MD5}")
	want, err := base64.StdEncoding.DecodeString(b64)
	if err != nil || len(want) != md5.Size {
		return false
	}
	sum := md5.Sum([]byte(password))
	return subtle.ConstantTimeCompare(sum[:], want) == 1
}
