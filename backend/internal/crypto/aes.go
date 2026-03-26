// Package crypto provides AES-256-GCM encryption/decryption for sensitive values
// (primarily LLM API keys stored in the database).
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"

	"github.com/skillora/backend/internal/config"
)

// Encrypt encrypts plaintext using AES-256-GCM with the AES_MASTER_KEY from config.
// The returned bytes are: [12-byte nonce] + [ciphertext + 16-byte GCM tag].
func Encrypt(plaintext string) ([]byte, error) {
	key, err := masterKey()
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("aes.NewCipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("cipher.NewGCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize()) // 12 bytes
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("nonce generation: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return ciphertext, nil
}

// Decrypt decrypts a ciphertext produced by Encrypt.
func Decrypt(ciphertext []byte) (string, error) {
	key, err := masterKey()
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("aes.NewCipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("cipher.NewGCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("crypto: ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("gcm.Open: %w", err)
	}

	return string(plaintext), nil
}

// masterKey decodes the 32-byte hex AES master key from config.
func masterKey() ([]byte, error) {
	raw := config.C.AESMasterKey
	key, err := hex.DecodeString(raw)
	if err != nil {
		return nil, fmt.Errorf("crypto: invalid AES_MASTER_KEY hex: %w", err)
	}
	if len(key) != 32 {
		return nil, fmt.Errorf("crypto: AES_MASTER_KEY must be 32 bytes (64 hex chars), got %d bytes", len(key))
	}
	return key, nil
}
