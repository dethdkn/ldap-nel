package passwords

import (
	"crypto/subtle"
	"encoding/base64"
	"strings"

	"golang.org/x/crypto/ripemd160"
)

func RMD160(password string) string {
	h := ripemd160.New()
	h.Write([]byte(password))
	sum := h.Sum(nil)
	return "{RMD160}" + base64.StdEncoding.EncodeToString(sum)
}

func VerifyRMD160(password, hashed string) bool {
	if !strings.HasPrefix(hashed, "{RMD160}") {
		return false
	}
	b64 := strings.TrimPrefix(hashed, "{RMD160}")
	want, err := base64.StdEncoding.DecodeString(b64)
	if err != nil || len(want) != ripemd160.Size {
		return false
	}

	h := ripemd160.New()
	h.Write([]byte(password))
	sum := h.Sum(nil)

	return subtle.ConstantTimeCompare(sum, want) == 1
}
