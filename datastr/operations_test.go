package datastr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	// Test Get method
	d := NewDict()

	// Test getting a value for a nonexistent key
	value := d.Get("nonexistent_key")
	assert.Nil(t, value, "Unexpected value for nonexistent key")

	// Test getting a value for an existing key
	d.Set("key1", "value1")
	value = d.Get("key1")
	assert.Equal(t, "value1", value, "Unexpected value for key1")
}

func TestSet(t *testing.T) {
	// // Test Set method
	d := NewDict()

	// Test setting a value for a nonexistent key
	d.Set("key1", "value1")
	value := d.Get("key1")
	assert.Equal(t, "value1", value, "Unexpected value for key1 after set")

	// Test updating a value for an existing key
	d.Set("key1", "updatedValue")
	value = d.Get("key1")
	assert.Equal(t, "updatedValue", value, "Unexpected value for key1 after update")
}

func TestDeleteMethod(t *testing.T) {
	// Test Delete method
	d := NewDict()

	// Test deleting a key that does not exist
	d.Delete("nonexistent_key")
	value := d.Get("nonexistent_key")
	assert.Nil(t, value, "Unexpected value for nonexistent key after delete")

	// Test deleting a key
	d.Set("key1", "value1")
	d.Delete("key1")
	value = d.Get("key1")
	assert.Nil(t, value, "Unexpected value for key1 after delete")
}
