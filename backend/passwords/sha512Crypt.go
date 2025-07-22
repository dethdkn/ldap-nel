package passwords

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/amoghe/go-crypt"
)

func generateRandomSalt(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base64.RawStdEncoding.EncodeToString(b)[:length]
}

// Generate a {CRYPT}$6$salt$... formatted hash
func Sha512Crypt(password, salt string) string {
	if len(salt) > 16 {
		panic("Salt must be 16 characters or less")
	}
	// If no salt is provided, generate a random one
	if salt == "" {
		salt = fmt.Sprintf("$6$%s$", generateRandomSalt(16))
	}
	hashed, err := crypt.Crypt(password, salt)
	if err != nil {
		panic(err)
	}
	return "{CRYPT}" + hashed
}

// Verifies a plaintext password against the stored {CRYPT} hash
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
