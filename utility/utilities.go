package utility

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

var errorBytes = `error generating random bytes`

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
			panic(errorBytes)
		}
	})
}

// GetRandomBytes returns a slice of bytes.
//
// No parameters.
// Returns [16]byte.
func GetRandomBytes() [16]byte {
	return randomBytes
}
