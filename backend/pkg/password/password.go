package password

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const (
	// Cost factor for bcrypt (12 as per security requirements)
	DefaultCost = 12
)

// Hash generates a bcrypt hash from a plain text password
func Hash(password string) (string, error) {
	if password == "" {
		return "", fmt.Errorf("password cannot be empty")
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hashedBytes), nil
}

// Verify checks if a plain text password matches a hashed password
func Verify(hashedPassword, password string) error {
	if hashedPassword == "" || password == "" {
		return fmt.Errorf("password and hash cannot be empty")
	}

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return fmt.Errorf("invalid password")
		}
		return fmt.Errorf("failed to verify password: %w", err)
	}

	return nil
}
