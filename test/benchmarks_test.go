package test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/dmarro89/go-redis-hashtable/datastr"

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

func BenchmarkSet(b *testing.B) {
	var d = datastr.NewDict()
	insertedElements := make(map[string]interface{})
	for i := 0; i < b.N; i++ {
		key := randomString(20)
		value := randomString(100)
		insertedElements[key] = value
	}
	b.ResetTimer()
	for key, value := range insertedElements {
		err := d.Set(key, value)
		b.StopTimer()
		assert.NoError(b, err, fmt.Sprintf("Error adding element {%s, %+v} to dictionary: %v", key, value, err))
		b.StartTimer()
	}
}

func BenchmarkGet(b *testing.B) {
	var d = datastr.NewDict()

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
		b.StopTimer()
		assert.NotNil(b, value, fmt.Sprintf("Error getting element %s from dictionary: %v", key, value))
		b.StartTimer()
	}
}

func BenchmarkDelete(b *testing.B) {
	var d = datastr.NewDict()

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
		b.StopTimer()
		assert.NoError(b, err, fmt.Sprintf("Error deleting element %s from dictionary: %v", key, err))
		b.StartTimer()
	}
}

func BenchmarkGoMapSet(b *testing.B) {
	insertedElements := make(map[string]interface{})
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		key := randomString(20)
		value := randomString(100)
		b.StartTimer()
		insertedElements[key] = value
	}
}

func BenchmarkGoMapGet(b *testing.B) {
	insertedElements := make(map[string]interface{})
	for i := 0; i < b.N; i++ {
		key := randomString(20)
		value := randomString(100)
		insertedElements[key] = value
	}
	b.ResetTimer()

	for key := range insertedElements {
		value := insertedElements[key]
		b.StopTimer()
		assert.NotNil(b, value, fmt.Sprintf("Error getting element %s from dictionary: %v", key, value))
		b.StartTimer()
	}
}

func BenchmarkGoMapDelete(b *testing.B) {
	insertedElements := make(map[string]interface{})
	keys := []string{}
	for i := 0; i < b.N; i++ {
		key := randomString(20)
		keys = append(keys, key)
		value := randomString(100)
		insertedElements[key] = value
	}
	b.ResetTimer()

	for _, value := range keys {
		delete(insertedElements, value)
		b.StopTimer()
		val := insertedElements[value]
		assert.Nil(b, val, fmt.Sprintf("Error deleting element %s from dictionary", value))
		b.StartTimer()
	}
}
