// Package hasher wraps stl package crypto by providing ability to implement specific
// cryptographic operations (creating hashes for passwords etc).
package hasher

import (
	"crypto/sha512"
	"fmt"
)

// HashPassword returns sha512 hash string for provided password
func HashPassword(password string) string {
	hash := sha512.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum(nil))
}

// CheckPasswordHash checks if provided password is equal to sha512 hash string
func CheckPasswordHash(password, hash string) bool {
	return HashPassword(password) == hash
}
