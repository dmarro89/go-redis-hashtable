package utilities

import (
	"crypto/rand"
	"sync"
)

var once sync.Once
var randomBytes [16]byte

// init initializes the program by generating 16 random bytes.
//
// No parameters.
// No return values.
func init() {
	generateRandomBytes()
}

// generateRandomBytes generates a slice of random bytes with the given length.
//
// Parameters:
// - length: an integer representing the length of the slice.
//
// Returns:
// - none
func generateRandomBytes() {
	once.Do(func() {
		_, err := rand.Read(randomBytes[:])
		if err != nil {
			panic(`error generating random bytes`)
		}
	})
}

// GetRandomBytes returns a slice of bytes.
//
// No parameters.
// Returns []byte.
func GetRandomBytes() [16]byte {
	return randomBytes
}
