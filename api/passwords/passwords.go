package passwords

import (
	"errors"
	"strings"
)

const cryptChars = "./0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

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
	case "Unix Crypt":
		pass, err = UnixCrypt(unencryptedPass, "")
	case "MD5 Crypt":
		pass, err = MD5Crypt(unencryptedPass, "")
	case "MD4":
		pass = MD4(unencryptedPass)
	case "MD5":
		pass = MD5(unencryptedPass)
	case "SHA1":
		pass = SHA1(unencryptedPass)
	case "SMD5":
		pass, err = SMD5(unencryptedPass, 8)
	case "SSHA":
		pass, err = SSHA(unencryptedPass, 8)
	case "SHA-256 Crypt":
		pass, err = Sha256Crypt(unencryptedPass, "")
	case "SHA-512 Crypt":
		pass, err = Sha512Crypt(unencryptedPass, "")
	case "RIPEMD-160":
		pass = RMD160(unencryptedPass)
	default:
		return "", errors.New("Unsupported encryption method: " + encryptionMethod)
	}
	if err != nil {
		return "", err
	}

	return pass, nil
}
