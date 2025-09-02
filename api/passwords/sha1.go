package passwords

import (
	"crypto/sha1"
	"crypto/subtle"
	"encoding/base64"
	"strings"
)

func SHA1(password string) string {
	sum := sha1.Sum([]byte(password))
	return "{SHA}" + base64.StdEncoding.EncodeToString(sum[:])
}

func VerifySHA1(password, hashed string) bool {
	if !strings.HasPrefix(hashed, "{SHA}") {
		return false
	}
	b64 := strings.TrimPrefix(hashed, "{SHA}")
	want, err := base64.StdEncoding.DecodeString(b64)
	if err != nil || len(want) != sha1.Size {
		return false
	}
	sum := sha1.Sum([]byte(password))
	return subtle.ConstantTimeCompare(sum[:], want) == 1
}
