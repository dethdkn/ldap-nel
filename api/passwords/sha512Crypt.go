package passwords

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/amoghe/go-crypt"
)

func saltSha512Crypt(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base64.RawStdEncoding.EncodeToString(b)[:length]
}

func Sha512Crypt(password, salt string) (string, error) {
	if len(salt) > 16 {
		return "", errors.New("salt must be 16 characters or less")
	}
	if salt == "" {
		salt = fmt.Sprintf("$6$%s$", saltSha512Crypt(16))
	}
	hashed, err := crypt.Crypt(password, salt)
	if err != nil {
		return "", err
	}
	return "{CRYPT}" + hashed, nil
}

func VerifySha512Crypt(password, hashed string) bool {
	if !strings.HasPrefix(hashed, "{CRYPT}") {
		return false
	}
	cryptHash := strings.TrimPrefix(hashed, "{CRYPT}")
	computed, err := crypt.Crypt(password, cryptHash)
	if err != nil {
		return false
	}
	return computed == cryptHash
}
