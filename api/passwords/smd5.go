package passwords

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"strings"
)

func SMD5(password string, saltLen int) (string, error) {
	if saltLen <= 0 {
		saltLen = 8
	}
	salt := make([]byte, saltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	h := md5.New()
	h.Write([]byte(password))
	h.Write(salt)
	digest := h.Sum(nil)

	final := append(digest, salt...)
	return "{SMD5}" + base64.StdEncoding.EncodeToString(final), nil
}

func VerifySMD5(password, hashed string) bool {
	if !strings.HasPrefix(hashed, "{SMD5}") {
		return false
	}
	b64 := strings.TrimPrefix(hashed, "{SMD5}")
	raw, err := base64.StdEncoding.DecodeString(b64)
	if err != nil || len(raw) <= md5.Size {
		return false
	}

	digest := raw[:md5.Size]
	salt := raw[md5.Size:]

	h := md5.New()
	h.Write([]byte(password))
	h.Write(salt)
	testDigest := h.Sum(nil)

	return subtle.ConstantTimeCompare(digest, testDigest) == 1
}
