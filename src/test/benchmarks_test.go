package test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/dmarro89/go-redis-hashtable/src/datastr"

	"github.com/stretchr/testify/assert"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randomString(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[r.Intn(len(letterBytes))]
	}
	return string(b)
}

var d = datastr.NewDict()

func BenchmarkSet(b *testing.B) {
	insertedElements := make(map[string]interface{})
	for i := 0; i < b.N; i++ {
		key := randomString(20)
		value := randomString(100)
		insertedElements[key] = value
	}
	b.ResetTimer()
	for key, value := range insertedElements {
		err := d.Set(key, value)
		assert.NoError(b, err, fmt.Sprintf("Error adding element {%s, %+v} to dictionary: %v", key, value, err))
	}
	b.ReportAllocs()
}

func BenchmarkGet(b *testing.B) {
	insertedElements := make(map[string]interface{})
	for i := 0; i < b.N; i++ {
		key := randomString(20)
		value := randomString(100)
		insertedElements[key] = value
		err := d.Set(key, value)
		assert.NoError(b, err, fmt.Sprintf("Error adding element {%s, %+v} to dictionary: %v", key, value, err))
	}
	b.ResetTimer()
	for key := range insertedElements {
		value := d.Get(key)
		assert.NotNil(b, value, fmt.Sprintf("Error getting element %s from dictionary: %v", key, value))
	}
	b.ReportAllocs()
}

func BenchmarkDelete(b *testing.B) {
	insertedElements := make(map[string]interface{})
	for i := 0; i < b.N; i++ {
		key := randomString(20)
		value := randomString(100)
		insertedElements[key] = value
		err := d.Set(key, value)
		assert.NoError(b, err, fmt.Sprintf("Error adding element {%s, %+v} to dictionary: %v", key, value, err))
	}
	b.ResetTimer()

	for key := range insertedElements {
		err := d.Delete(key)
		assert.NoError(b, err, fmt.Sprintf("Error deleting element %s from dictionary: %v", key, err))
	}
	b.ReportAllocs()
}
