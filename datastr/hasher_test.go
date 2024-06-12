package datastr

import (
	"testing"

	"github.com/dchest/siphash"
	"github.com/stretchr/testify/assert"
)

func TestNewHasher(t *testing.T) {
	hasher := NewHasher()
	assert.NotNil(t, hasher, "NewHasher() should not return nil")

	assert.NotEqual(t, uint64(0), hasher.key0, "NewHasher() should set key0")
	assert.NotEqual(t, uint64(0), hasher.key1, "NewHasher() should set key1")
}

func TestDigest(t *testing.T) {
	hasher := NewHasher()
	message := "test message"
	expectedHash := siphash.Hash(hasher.key0, hasher.key1, []byte(message))

	actualHash := hasher.Digest(message)
	assert.Equal(t, expectedHash, actualHash, "Digest(%q) should return correct hash", message)
}

func TestDigestBufferReused(t *testing.T) {
	hasher := NewHasher()
	message := "test message"

	hasher.Digest(message)
	initialBuf := customPool.Get().(*[]byte)
	initialCap := cap(*initialBuf)
	customPool.Put(initialBuf)

	hasher.Digest(message)
	secondBuf := customPool.Get().(*[]byte)
	secondCap := cap(*secondBuf)
	customPool.Put(secondBuf)

	assert.Equal(t, initialCap, secondCap, "Buffer capacity should be reused correctly")
}
