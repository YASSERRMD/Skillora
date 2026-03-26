package crypto_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/skillora/backend/internal/config"
	"github.com/skillora/backend/internal/crypto"
)

func init() {
	os.Setenv("DATABASE_URL", "postgres://test:test@localhost:5432/testdb")
	os.Setenv("REDIS_URL", "redis://localhost:6379")
	os.Setenv("GOOGLE_CLIENT_ID", "test-client-id")
	os.Setenv("GOOGLE_CLIENT_SECRET", "test-client-secret")
	os.Setenv("JWT_SECRET", "test-jwt-secret")
	// Valid 32-byte (64 hex char) key.
	os.Setenv("AES_MASTER_KEY", "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2")
	config.Load()
}

func TestEncryptDecryptRoundTrip(t *testing.T) {
	plaintext := "sk-secretOpenAIKey-abc123"
	ciphertext, err := crypto.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt error: %v", err)
	}
	if len(ciphertext) == 0 {
		t.Fatal("Encrypt returned empty ciphertext")
	}

	got, err := crypto.Decrypt(ciphertext)
	if err != nil {
		t.Fatalf("Decrypt error: %v", err)
	}
	if got != plaintext {
		t.Errorf("round-trip mismatch: got %q, want %q", got, plaintext)
	}
}

func TestEncryptProducesUniqueCiphertexts(t *testing.T) {
	// Two encryptions of the same plaintext must differ (random nonce).
	plain := "same-key"
	ct1, _ := crypto.Encrypt(plain)
	ct2, _ := crypto.Encrypt(plain)
	if bytes.Equal(ct1, ct2) {
		t.Error("two encryptions of the same plaintext should differ (random nonce)")
	}
}

func TestDecryptTamperedCiphertext(t *testing.T) {
	ct, _ := crypto.Encrypt("tamper-me")
	ct[len(ct)-1] ^= 0xFF // flip last byte
	if _, err := crypto.Decrypt(ct); err == nil {
		t.Error("expected error when decrypting tampered ciphertext")
	}
}

func TestDecryptTooShort(t *testing.T) {
	if _, err := crypto.Decrypt([]byte{1, 2, 3}); err == nil {
		t.Error("expected error for too-short ciphertext")
	}
}

func TestEncryptEmptyString(t *testing.T) {
	ct, err := crypto.Encrypt("")
	if err != nil {
		t.Fatalf("Encrypt empty: %v", err)
	}
	got, err := crypto.Decrypt(ct)
	if err != nil {
		t.Fatalf("Decrypt empty round-trip: %v", err)
	}
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}
