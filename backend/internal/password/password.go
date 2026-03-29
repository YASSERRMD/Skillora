package password

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrPasswordTooShort is returned when password is less than 8 characters
	ErrPasswordTooShort = errors.New("password must be at least 8 characters long")

	// ErrPasswordMismatch is returned when password doesn't match hash
	ErrPasswordMismatch = errors.New("password does not match")
)

const (
	// MinPasswordLength is the minimum allowed password length
	MinPasswordLength = 8

	// BcryptCost is the cost factor for bcrypt hashing
	BcryptCost = 12
)

// HashPassword creates a bcrypt hash from a plain-text password
func HashPassword(password string) (string, error) {
	if len(password) < MinPasswordLength {
		return "", fmt.Errorf("%w: minimum %d characters required", ErrPasswordTooShort, MinPasswordLength)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), BcryptCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hash), nil
}

// VerifyPassword checks if a plain-text password matches a bcrypt hash
func VerifyPassword(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return ErrPasswordMismatch
	}
	return nil
}

// ValidatePassword checks if a password meets security requirements
func ValidatePassword(password string) error {
	if len(password) < MinPasswordLength {
		return fmt.Errorf("%w: minimum %d characters required", ErrPasswordTooShort, MinPasswordLength)
	}
	return nil
}
