package passwords

import (
	"crypto/rand"
	"crypto/subtle"
	"errors"
	"strings"

	"github.com/amoghe/go-crypt"
)

func saltSha256Crypt(n int) (string, error) {
	if n < 1 || n > 16 {
		return "", errors.New("salt length must be between 1 and 16 for SHA-256 crypt")
	}
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	s := make([]byte, n)
	for i, v := range b {
		s[i] = cryptChars[int(v)%len(cryptChars)]
	}
	return string(s), nil
}

func Sha256Crypt(password, salt string) (string, error) {
	if salt == "" {
		var err error
		salt, err = saltSha256Crypt(16)
		if err != nil {
			return "", err
		}
	}
	if len(salt) > 16 {
		return "", errors.New("salt must be 16 characters or less for SHA-256 crypt")
	}
	fullSalt := "$5$" + salt
	h, err := crypt.Crypt(password, fullSalt)
	if err != nil {
		return "", err
	}
	return "{CRYPT}" + h, nil
}

func VerifySha256Crypt(password, hashed string) bool {
	if !strings.HasPrefix(hashed, "{CRYPT}$5$") {
		return false
	}
	cryptHash := strings.TrimPrefix(hashed, "{CRYPT}")

	parts := strings.SplitN(cryptHash, "$", 4)
	if len(parts) < 4 {
		return false
	}
	salt := parts[2]
	fullSalt := "$5$" + salt

	computed, err := crypt.Crypt(password, fullSalt)
	if err != nil {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(computed), []byte(cryptHash)) == 1
}
