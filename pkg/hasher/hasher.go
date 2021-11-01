// Package hasher wraps stl package crypto by providing ability to implement specific
// cryptographic operations (creating hashes for passwords etc).
package hasher

import (
	"crypto/sha512"
	"fmt"
)

// HashPassword returns sha512 hash string for provided password
func HashPassword(password string) (string, error) {
	hash := sha512.New()
	if _, err := hash.Write([]byte(password)); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// CheckPasswordHash checks if provided password is equal to sha512 hash string
// In case there was error when composing password hash false value will be returned
func CheckPasswordHash(password, hash string) bool {
	passwordHash, err := HashPassword(password)
	if err != nil {
		return false
	}

	return passwordHash == hash
}
