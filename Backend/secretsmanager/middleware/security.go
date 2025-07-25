package middleware

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"secretsmanager/constants"

	siv "github.com/secure-io/siv-go"
)

var key = []byte(constants.AES256_32BYTE_ENCRYPTION_KEY)
var aesSIVKey = []byte(constants.AES256_32BYTE_DETERMINISTIC_ENCRYPTION_KEY) // 32 or 64 bytes

// Encrypt encrypts plain text using AES-GCM
func Encrypt(plainText string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	cipherText := aesGCM.Seal(nonce, nonce, []byte(plainText), nil)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// EncryptDeterministic encrypts using AES-SIV (deterministic)
func EncryptDeterministic(plainText string) (string, error) {
	sivCipher, err := siv.NewCMAC(aesSIVKey)
	if err != nil {
		return "", err
	}

	cipherText := sivCipher.Seal(nil, nil, []byte(plainText), nil)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// Decrypt decrypts cipher text using AES-GCM
func Decrypt(encrypted string) (string, error) {
	cipherText, err := base64.StdEncoding.DecodeString(encrypted)
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
	if len(cipherText) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]
	plainText, err := aesGCM.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}

// DecryptDeterministic decrypts AES-SIV cipher text
func DecryptDeterministic(cipherTextB64 string) (string, error) {
	sivCipher, err := siv.NewCMAC(aesSIVKey)
	if err != nil {
		return "", err
	}

	cipherText, err := base64.StdEncoding.DecodeString(cipherTextB64)
	if err != nil {
		return "", err
	}

	plainText, err := sivCipher.Open(nil, nil, cipherText, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}
