package utilities

import (
	"crypto/rand"
	"sync"
)

var once sync.Once
var randomBytes []byte

// init initializes the program by generating 16 random bytes.
//
// No parameters.
// No return values.
func init() {
	generateRandomBytes(16)
}

// generateRandomBytes generates a slice of random bytes with the given length.
//
// Parameters:
// - length: an integer representing the length of the slice.
//
// Returns:
// - none
func generateRandomBytes(length int) {
	once.Do(func() {
		randomBytes = make([]byte, length)
		_, err := rand.Read(randomBytes)
		if err != nil {
			panic(`error generating random bytes`)
		}
	})
}

// GetRandomBytes returns a slice of bytes.
//
// No parameters.
// Returns []byte.
func GetRandomBytes() []byte {
	return randomBytes
}
