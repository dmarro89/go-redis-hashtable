package structure

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDictEntry(t *testing.T) {
	key := "testKey"
	value := "testValue"

	entry := NewDictEntry(key, value)
	assert.Equal(t, key, entry.key, "Expected key %s, but got %s", key, entry.key)
	assert.Equal(t, value, entry.value, "Expected value %v, but got %v", value, entry.value)
	assert.Nil(t, entry.next, "Expected next to be nil, but it's not")
}

func TestDictEntryNext(t *testing.T) {
	entry1 := NewDictEntry("key1", "value1")
	entry2 := NewDictEntry("key2", "value2")

	entry1.next = entry2
	assert.Equal(t, entry2, entry1.next, "Expected entry1.next to be entry2, but it's not")
}
