package test

import (
	"fmt"
	"testing"

	"github.com/dmarro89/go-redis-hashtable/src/datastr"

	"github.com/stretchr/testify/assert"
)

func TestSequentialOperations(t *testing.T) {
	d := datastr.NewDict()

	const numberOfOperations = 10000

	insertedElements := make(map[string]interface{})
	for i := 0; i < numberOfOperations; i++ {
		key := randomString(20)
		value := randomString(30)
		err := d.Set(key, value)
		assert.NoError(t, err)
	}

	for key, expectedValue := range insertedElements {
		assert.Equal(t, expectedValue, d.Get(key))
	}

	// Delete the inserted elements
	for key := range insertedElements {
		err := d.Delete(key)
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
}
