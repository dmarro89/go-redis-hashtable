package hashing

import (
	"testing"

	"github.com/dchest/siphash"
	"github.com/stretchr/testify/assert"
)

func TestNewHasher(t *testing.T) {
	hasher := NewSip24Hasher().(*Sip24Hasher)
	assert.NotNil(t, hasher, "NewHasher() should not return nil")

	assert.NotEqual(t, uint64(0), hasher.Key0, "NewHasher() should set key0")
	assert.NotEqual(t, uint64(0), hasher.Key1, "NewHasher() should set key1")
}

func TestDigest(t *testing.T) {
	hasher := NewSip24Hasher().(*Sip24Hasher)
	message := "test message"
	expectedHash := siphash.Hash(hasher.Key0, hasher.Key1, []byte(message))

	actualHash := hasher.Digest(message)
	assert.Equal(t, expectedHash, actualHash, "Digest(%q) should return correct hash", message)
}
