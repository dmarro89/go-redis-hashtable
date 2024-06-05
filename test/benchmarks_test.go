package test

import (
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/dmarro89/go-redis-hashtable/datastr"

	"github.com/stretchr/testify/assert"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const maxStringLength = 100

func randomString(length int) string {
	if length == 0 {
		length = rand.IntN(maxStringLength)
	}

	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[rand.IntN(len(letterBytes))]
	}
	return string(b)
}

var Set100 = make(map[string]interface{})
var Set1000 = make(map[string]interface{})
var Set3000 = make(map[string]interface{})
var Set10000 = make(map[string]interface{})
var dict100 = datastr.NewDict()
var dict1000 = datastr.NewDict()
var dict3000 = datastr.NewDict()
var dict10000 = datastr.NewDict()
var map100 = make(map[string]interface{})
var map1000 = make(map[string]interface{})
var map3000 = make(map[string]interface{})
var map10000 = make(map[string]interface{})

func init() {
	for i := 0; i < 100; i++ {
		Set100[randomString(0)] = randomString(0)
	}
	for i := 0; i < 1000; i++ {
		Set1000[randomString(0)] = randomString(0)
	}
	for i := 0; i < 3000; i++ {
		Set3000[randomString(0)] = randomString(0)
	}
	for i := 0; i < 10000; i++ {
		Set10000[randomString(0)] = randomString(0)
	}
}

func BenchmarkSet100(b *testing.B) {
	b.ResetTimer()
	for key, value := range Set100 {
		b.StartTimer()
		err := dict100.Set(key, value)
		b.StopTimer()
		assert.NoError(b, err, fmt.Sprintf("Error adding element {%s, %+v} to dictionary: %v", key, value, err))
	}
}

func BenchmarkSet1000(b *testing.B) {
	b.ResetTimer()
	for key, value := range Set1000 {
		b.StartTimer()
		err := dict1000.Set(key, value)
		b.StopTimer()
		assert.NoError(b, err, fmt.Sprintf("Error adding element {%s, %+v} to dictionary: %v", key, value, err))
	}
}

func BenchmarkSet3000(b *testing.B) {
	b.ResetTimer()
	for key, value := range Set3000 {
		b.StartTimer()
		err := dict3000.Set(key, value)
		b.StopTimer()
		assert.NoError(b, err, fmt.Sprintf("Error adding element {%s, %+v} to dictionary: %v", key, value, err))
	}
}

func BenchmarkSet10000(b *testing.B) {
	b.ResetTimer()
	for key, value := range Set10000 {
		b.StartTimer()
		err := dict10000.Set(key, value)
		b.StopTimer()
		assert.NoError(b, err, fmt.Sprintf("Error adding element {%s, %+v} to dictionary: %v", key, value, err))
	}
}

func BenchmarkGet100(b *testing.B) {
	for key, val := range Set100 {
		b.StartTimer()
		value := dict100.Get(key)
		b.StopTimer()
		assert.Equal(b, value, val, fmt.Sprintf("Error getting element %s from dictionary: %v", key, value))
	}
}

func BenchmarkGet1000(b *testing.B) {
	for key, val := range Set1000 {
		b.StartTimer()
		value := dict1000.Get(key)
		b.StopTimer()
		assert.Equal(b, value, val, fmt.Sprintf("Error getting element %s from dictionary: %v", key, value))
	}
}

func BenchmarkGet3000(b *testing.B) {
	for key, val := range Set3000 {
		b.StartTimer()
		value := dict3000.Get(key)
		b.StopTimer()
		assert.Equal(b, value, val, fmt.Sprintf("Error getting element %s from dictionary: %v", key, value))
	}
}

func BenchmarkGet10000(b *testing.B) {
	for key, val := range Set10000 {
		b.StartTimer()
		value := dict10000.Get(key)
		b.StopTimer()
		assert.Equal(b, value, val, fmt.Sprintf("Error getting element %s from dictionary: %v", key, value))
	}
}

func BenchmarkDelete100(b *testing.B) {
	for key := range Set100 {
		b.StartTimer()
		err := dict100.Delete(key)
		b.StopTimer()
		assert.NoError(b, err, fmt.Sprintf("Error deleting element %s from dictionary: %v", key, err))
	}
}

func BenchmarkDelete1000(b *testing.B) {
	for key := range Set1000 {
		b.StartTimer()
		err := dict1000.Delete(key)
		b.StopTimer()
		assert.NoError(b, err, fmt.Sprintf("Error deleting element %s from dictionary: %v", key, err))
	}
}

func BenchmarkDelete3000(b *testing.B) {
	for key := range Set3000 {
		b.StartTimer()
		err := dict3000.Delete(key)
		b.StopTimer()
		assert.NoError(b, err, fmt.Sprintf("Error deleting element %s from dictionary: %v", key, err))
	}
}

func BenchmarkDelete10000(b *testing.B) {
	for key := range Set10000 {
		b.StartTimer()
		err := dict10000.Delete(key)
		b.StopTimer()
		assert.NoError(b, err, fmt.Sprintf("Error deleting element %s from dictionary: %v", key, err))
	}
}

func BenchmarkGoMapSet100(b *testing.B) {
	b.ResetTimer()
	for key, value := range Set100 {
		b.StartTimer()
		map100[key] = value
		b.StopTimer()
	}
}

func BenchmarkGoMapSet1000(b *testing.B) {
	b.ResetTimer()
	for key, value := range Set1000 {
		b.StartTimer()
		map1000[key] = value
		b.StopTimer()
	}
}

func BenchmarkGoMapSet3000(b *testing.B) {
	b.ResetTimer()
	for key, value := range Set3000 {
		b.StartTimer()
		map3000[key] = value
		b.StopTimer()
	}
}

func BenchmarkGoMapSet10000(b *testing.B) {
	b.ResetTimer()
	for key, value := range Set10000 {
		b.StartTimer()
		map10000[key] = value
		b.StopTimer()
	}
}

func BenchmarkGoMapGet100(b *testing.B) {
	for key, val := range Set100 {
		b.StartTimer()
		value := map100[key]
		b.StopTimer()
		assert.Equal(b, value, val, fmt.Sprintf("Error getting element %s from dictionary: %v", key, value))
	}
}

func BenchmarkGoMapGet1000(b *testing.B) {
	for key, val := range Set1000 {
		b.StartTimer()
		value := map1000[key]
		b.StopTimer()
		assert.Equal(b, value, val, fmt.Sprintf("Error getting element %s from dictionary: %v", key, value))
	}
}

func BenchmarkGoMapGet3000(b *testing.B) {
	for key, val := range Set3000 {
		b.StartTimer()
		value := map3000[key]
		b.StopTimer()
		assert.Equal(b, value, val, fmt.Sprintf("Error getting element %s from dictionary: %v", key, value))
	}
}

func BenchmarkGoMapGet10000(b *testing.B) {
	for key, val := range Set10000 {
		b.StartTimer()
		value := map10000[key]
		b.StopTimer()
		assert.Equal(b, value, val, fmt.Sprintf("Error getting element %s from dictionary: %v", key, value))
	}
}

func BenchmarkGoMapDelete100(b *testing.B) {
	for key := range Set100 {
		b.StartTimer()
		delete(map100, key)
		b.StopTimer()
	}
}

func BenchmarkGoMapDelete1000(b *testing.B) {
	for key := range Set1000 {
		b.StartTimer()
		delete(map1000, key)
		b.StopTimer()
	}
}

func BenchmarkGoMapDelete3000(b *testing.B) {
	for key := range Set3000 {
		b.StartTimer()
		delete(map3000, key)
		b.StopTimer()
	}
}

func BenchmarkGoMapDelete10000(b *testing.B) {
	for key := range Set10000 {
		b.StartTimer()
		delete(map10000, key)
		b.StopTimer()
	}
}
