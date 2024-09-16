// secure_test.go
package secure

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"io"
	"testing"
)

func TestEncryptDecryptAES(t *testing.T) {
	key := []byte("a very strong encryption key 123") // Must be 16, 24, or 32 bytes long
	plainText := "This is a secret message!"

	// Test encryption
	cipherText, err := EncryptAES(plainText, key)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	// Ensure that cipherText is not empty
	if cipherText == "" {
		t.Fatal("cipherText is empty")
	}

	// Test decryption
	decryptedText, err := DecryptAES(cipherText, key)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	// Ensure the decrypted text matches the original plaintext
	if decryptedText != plainText {
		t.Fatalf("Decryption produced incorrect result: got %v, want %v", decryptedText, plainText)
	}
}

func TestDecryptAESWithInvalidKey(t *testing.T) {
	key := []byte("a very strong encryption key 123")
	wrongKey := []byte("wrong key! 1234567890")
	plainText := "This is another secret message!"

	// Test encryption
	cipherText, err := EncryptAES(plainText, key)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	// Try to decrypt with the wrong key
	_, err = DecryptAES(cipherText, wrongKey)
	if err == nil {
		t.Fatal("Expected decryption to fail with the wrong key, but it succeeded")
	}
}

func TestDecryptAESWithInvalidCipherText(t *testing.T) {
	key := []byte("a very strong encryption key 123")

	// Test decryption with an invalid (short) cipher text
	invalidCipherText := "abcd" // Shorter than AES block size
	_, err := DecryptAES(invalidCipherText, key)
	if err != ErrCipherTextTooShort {
		t.Fatalf("Expected ErrCipherTextTooShort, but got %v", err)
	}
}

func TestDeriveKeyPBKDF2(t *testing.T) {
	password := []byte("securepassword")
	salt := []byte("randomsalt")

	key := DeriveKeyPBKDF2(password, salt)

	// Ensure that the key has the correct length
	expectedKeyLen := 32
	if len(key) != expectedKeyLen {
		t.Fatalf("Expected key length of %d, but got %d", expectedKeyLen, len(key))
	}

	// Deriving the key with the same password and salt should produce the same result
	key2 := DeriveKeyPBKDF2(password, salt)
	if !bytes.Equal(key, key2) {
		t.Fatal("Deriving the same key twice did not produce the same result")
	}
}

func TestDeriveKeyScrypt(t *testing.T) {
	password := []byte("securepassword")
	salt := []byte("randomsalt")

	key, err := DeriveKeyScrypt(password, salt)
	if err != nil {
		t.Fatalf("Scrypt key derivation failed: %v", err)
	}

	// Ensure that the key has the correct length
	expectedKeyLen := 32
	if len(key) != expectedKeyLen {
		t.Fatalf("Expected key length of %d, but got %d", expectedKeyLen, len(key))
	}

	// Deriving the key with the same password and salt should produce the same result
	key2, err := DeriveKeyScrypt(password, salt)
	if err != nil {
		t.Fatalf("Scrypt key derivation failed: %v", err)
	}
	if !bytes.Equal(key, key2) {
		t.Fatal("Deriving the same key twice did not produce the same result")
	}
}

func TestHexEncodeDecode(t *testing.T) {
	original := "Hello, World!"
	encoded := hex.EncodeToString([]byte(original))

	decodedBytes, err := hex.DecodeString(encoded)
	if err != nil {
		t.Fatalf("Failed to decode hex: %v", err)
	}
	decoded := string(decodedBytes)

	if decoded != original {
		t.Fatalf("Expected %v, but got %v", original, decoded)
	}
}

// test password verification

// Test to check if decryption with the correct password works
func TestMasterPasswordValidation(t *testing.T) {
	// Given known plaintext and salt
	plainText := "this is a secret"
	password := []byte("correct_password")
	wrongPassword := []byte("wrong_password")
	salt := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		t.Fatalf("Failed to generate salt: %v", err)
	}

	// Derive key from the correct password using PBKDF2
	key := DeriveKeyPBKDF2(password, salt)

	// Encrypt the plaintext with the correct key
	encrypted, err := EncryptAES(plainText, key)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	// Test decryption with the correct password
	derivedKey := DeriveKeyPBKDF2(password, salt)
	decryptedText, err := DecryptAES(encrypted, derivedKey)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	if decryptedText != plainText {
		t.Fatalf("Decrypted text doesn't match the original")
	}

	// Test decryption with the wrong password
	wrongDerivedKey := DeriveKeyPBKDF2(wrongPassword, salt)
	wrongM, err2 := DecryptAES(encrypted, wrongDerivedKey)
	if err2 != nil {
		t.Fatalf("Decryption failed")
	} else {
		if wrongM == plainText {
			t.Fatalf("Decryption with wrong password worked! We expected '%v' not to be equal to '%v'", wrongM, plainText)
		}
	}
}

// Utility to create hex-encoded ciphertext from string
func hexDecode(cipherHex string) []byte {
	cipher, _ := hex.DecodeString(cipherHex)
	return cipher
}

// Utility to simulate derived key mismatch case (if needed for further validation)
func TestMasterPasswordValidationMismatch(t *testing.T) {
	plainText := "super_secret"
	password := []byte("strong_password")
	salt := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		t.Fatalf("Failed to generate salt: %v", err)
	}

	// Encrypt using the correct password
	key, _ := DeriveKeyScrypt(password, salt)
	cipherText, err := EncryptAES(plainText, key)
	if err != nil {
		t.Fatalf("Encryption error: %v", err)
	}

	// Try to decrypt with wrong password
	wrongPassword := []byte("wrong_password")
	wrongKey, _ := DeriveKeyScrypt(wrongPassword, salt)
	decryptedText, err := DecryptAES(cipherText, wrongKey)
	if err != nil {
		t.Fatal("Decryption failed")
	}

	// Verify wrong decryption output should not match plaintext
	if bytes.Equal([]byte(decryptedText), []byte(plainText)) {
		t.Fatal("Mismatch expected, but decryption returned original text")
	}
}
