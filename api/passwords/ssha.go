package passwords

import (
	"crypto/rand"
	"crypto/sha1"
	"crypto/subtle"
	"encoding/base64"
	"strings"
)

func SSHA(password string, saltLen int) (string, error) {
	if saltLen <= 0 {
		saltLen = 8
	}
	salt := make([]byte, saltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	h := sha1.New()
	h.Write([]byte(password))
	h.Write(salt)
	digest := h.Sum(nil)

	final := append(digest, salt...)
	return "{SSHA}" + base64.StdEncoding.EncodeToString(final), nil
}

func VerifySSHA(password, hashed string) bool {
	if !strings.HasPrefix(hashed, "{SSHA}") {
		return false
	}
	b64 := strings.TrimPrefix(hashed, "{SSHA}")
	raw, err := base64.StdEncoding.DecodeString(b64)
	if err != nil || len(raw) <= sha1.Size {
		return false
	}

	digest := raw[:sha1.Size]
	salt := raw[sha1.Size:]

	h := sha1.New()
	h.Write([]byte(password))
	h.Write(salt)
	testDigest := h.Sum(nil)

	return subtle.ConstantTimeCompare(digest, testDigest) == 1
}
