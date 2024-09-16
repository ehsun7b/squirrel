// secure.go
package secure

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"

	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/scrypt"
)

// Custom error for ciphertext that is too short.
var ErrCipherTextTooShort = errors.New("ciphertext too short")

// EncryptAES encrypts plaintext using AES encryption.
func EncryptAES(plainText string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	plainTextBytes := []byte(plainText)
	cipherText := make([]byte, aes.BlockSize+len(plainTextBytes))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainTextBytes)

	return hex.EncodeToString(cipherText), nil
}

// DecryptAES decrypts AES-encrypted ciphertext.
func DecryptAES(cipherTextHex string, key []byte) (string, error) {
	cipherText, err := hex.DecodeString(cipherTextHex)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(cipherText) < aes.BlockSize {
		return "", ErrCipherTextTooShort
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), nil
}

// DeriveKeyPBKDF2 derives a key from a password using PBKDF2.
func DeriveKeyPBKDF2(password, salt []byte) []byte {
	return pbkdf2.Key(password, salt, 4096, 32, sha256.New)
}

// DeriveKeyScrypt derives a key from a password using scrypt.
func DeriveKeyScrypt(password, salt []byte) ([]byte, error) {
	return scrypt.Key(password, salt, 16384, 8, 1, 32)
}

func GenerateSalt(size int) ([]byte, error) {
	salt := make([]byte, size)
	// Fill salt with random bytes
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}
