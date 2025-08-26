package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"os"
)

func Encrypt(text string) (string, error) {
	encryptionKey := os.Getenv("ENCRYPTION_KEY")
	if encryptionKey == "" {
		return "", errors.New("ENCRYPTION_KEY is not set")
	}

	key, err := base64.StdEncoding.DecodeString(encryptionKey)
	if err != nil {
		return "", errors.New("ENCRYPTION_KEY must be valid base64")
	}

	if len(key) != 32 {
		return "", errors.New("ENCRYPTION_KEY must be exactly 32 bytes after base64 decoding")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, []byte(text), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(encrypted string) (string, error) {
	encryptionKey := os.Getenv("ENCRYPTION_KEY")
	if encryptionKey == "" {
		return "", errors.New("ENCRYPTION_KEY is not set")
	}

	key, err := base64.StdEncoding.DecodeString(encryptionKey)
	if err != nil {
		return "", errors.New("ENCRYPTION_KEY must be valid base64")
	}

	if len(key) != 32 {
		return "", errors.New("ENCRYPTION_KEY must be exactly 32 bytes after base64 decoding")
	}

	data, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("malformed ciphertext")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
