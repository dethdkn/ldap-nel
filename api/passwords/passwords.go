package passwords

import (
	"errors"
	"strings"
)

func GetEncryptedPass(value string) (string, error) {
	parts := strings.SplitN(value, ":", 2)
	if len(parts) != 2 {
		return "", errors.New("Invalid password format.")
	}

	encryptionMethod := strings.TrimSpace(parts[0])
	unencryptedPass := strings.TrimSpace(parts[1])

	var pass string
	var err error

	switch encryptionMethod {
	case "Plain Text":
		pass = unencryptedPass
	case "SHA-512 Crypt":
		pass, err = Sha512Crypt(unencryptedPass, "")
	}
	if err != nil {
		return "", err
	}

	return pass, nil
}
