package structure

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHashTable(t *testing.T) {
	size := int64(10)
	ht := NewHashTable(size)

	assert.Equal(t, size, ht.size, "Expected size %d, but got %d", size, ht.size)
	assert.Equal(t, uint64(size-1), ht.sizemask, "Expected sizemask %d, but got %d", size-1, ht.sizemask)
	assert.NotNil(t, ht.table, "Expected table to be initialized, but it's nil")
	assert.False(t, ht.empty(), "Expected HashTable to be not empty, but it's empty")
}

func TestEmptyHashTable(t *testing.T) {
	size := int64(0)
	ht := NewHashTable(0)

	assert.Equal(t, size, ht.size, "Expected size %d, but got %d", size, ht.size)
	assert.Equal(t, uint64(size), ht.sizemask, "Expected sizemask %d, but got %d", size, ht.sizemask)
	assert.Empty(t, ht.table, "Expected table to be nil but it's initialized")
	assert.True(t, ht.empty(), "Expected HashTable to be empty, but it's not")
}
