package passwords

import (
	"crypto/rand"
	"errors"
	"strings"

	"github.com/amoghe/go-crypt"
)

func saltUnixCrypt(n int) (string, error) {
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

func UnixCrypt(password, salt string) (string, error) {
	if salt == "" {
		var err error
		salt, err = saltUnixCrypt(2)
		if err != nil {
			return "", err
		}
	}
	if len(salt) != 2 {
		return "", errors.New("salt must be exactly 2 characters for DES crypt")
	}
	h, err := crypt.Crypt(password, salt)
	if err != nil {
		return "", err
	}
	return "{CRYPT}" + h, nil
}

func VerifyUnixCrypt(password, hashed string) bool {
	if !strings.HasPrefix(hashed, "{CRYPT}") {
		return false
	}
	cryptHash := strings.TrimPrefix(hashed, "{CRYPT}")
	if len(cryptHash) != 13 {
		return false
	}
	computed, err := crypt.Crypt(password, cryptHash)
	if err != nil {
		return false
	}
	return computed == cryptHash
}
