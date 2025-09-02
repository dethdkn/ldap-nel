package passwords

import (
	"crypto/subtle"
	"encoding/base64"
	"hash"
	"strings"

	"golang.org/x/crypto/md4"
)

func md4Sum(data []byte) []byte {
	var h hash.Hash = md4.New()
	h.Write(data)
	return h.Sum(nil)
}

func MD4(password string) string {
	sum := md4Sum([]byte(password))
	return "{MD4}" + base64.StdEncoding.EncodeToString(sum)
}

func VerifyMD4(password, hashed string) bool {
	if !strings.HasPrefix(hashed, "{MD4}") {
		return false
	}
	b64 := strings.TrimPrefix(hashed, "{MD4}")
	want, err := base64.StdEncoding.DecodeString(b64)
	if err != nil || len(want) != md4.Size {
		return false
	}
	sum := md4Sum([]byte(password))
	return subtle.ConstantTimeCompare(sum, want) == 1
}
