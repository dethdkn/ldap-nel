package passwords

import (
	"crypto/rand"
	"crypto/subtle"
	"errors"
	"strings"

	"github.com/amoghe/go-crypt"
)

func saltMD5Crypt(n int) (string, error) {
	if n < 1 || n > 16 {
		return "", errors.New("salt length must be between 1 and 16 for MD5 crypt")
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

func MD5Crypt(password, salt string) (string, error) {
	if salt == "" {
		var err error
		salt, err = saltMD5Crypt(8)
		if err != nil {
			return "", err
		}
	}
	if len(salt) > 16 {
		return "", errors.New("salt must be 16 characters or less for MD5 crypt")
	}
	fullSalt := "$1$" + salt
	h, err := crypt.Crypt(password, fullSalt)
	if err != nil {
		return "", err
	}
	return "{CRYPT}" + h, nil
}

func VerifyMD5Crypt(password, hashed string) bool {
	if !strings.HasPrefix(hashed, "{CRYPT}$1$") {
		return false
	}
	cryptHash := strings.TrimPrefix(hashed, "{CRYPT}")
	parts := strings.SplitN(cryptHash, "$", 4)
	if len(parts) < 4 {
		return false
	}
	salt := parts[2]
	fullSalt := "$1$" + salt
	computed, err := crypt.Crypt(password, fullSalt)
	if err != nil {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(computed), []byte(cryptHash)) == 1
}
