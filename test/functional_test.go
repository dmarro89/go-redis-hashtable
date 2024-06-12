package test

import (
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/dmarro89/go-redis-hashtable/datastr"

	"github.com/stretchr/testify/assert"
)

const maxStringLength = 100
const letterBytes = "abcdefghijklmnopqrstuvwxyz"

func randomString(length int) string {
	if length == 0 {
		length = rand.IntN(maxStringLength) + 1
	}

	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[rand.IntN(len(letterBytes))]
	}
	return string(b)
}

func TestSequentialOperations(t *testing.T) {
	d := datastr.NewDict()

	const numberOfOperations = 1000000

	insertedElements := make(map[string]interface{})
	for i := 0; i < numberOfOperations; i++ {
		key := randomString(0)
		value := randomString(0)
		insertedElements[key] = value
		err := d.Set(key, value)
		assert.NoError(t, err)
	}

	// Attempt to delete non-existent elements
	for i := 0; i < numberOfOperations; i++ {
		key := fmt.Sprintf("nonexistent%d", i)
		err := d.Delete(key)
		assert.EqualError(t, err, `entry not found`)
	}

	// Verify the retrieval of N non-existent elements
	for i := 0; i < numberOfOperations; i++ {
		key := fmt.Sprintf("nonexistent%d", i)
		assert.Nil(t, d.Get(key))
	}

	for key, expectedValue := range insertedElements {
		assert.Equal(t, expectedValue, d.Get(key))
	}

	// Delete the inserted elements
	for key := range insertedElements {
		err := d.Delete(key)
		assert.NoError(t, err)
	}
}
