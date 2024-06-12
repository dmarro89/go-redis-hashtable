package test

import (
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/dmarro89/go-redis-hashtable/structure"
)

type keyValue struct{ Key, Value string }

const letterBytes = "abcdefghijklmnopqrstuvwxyz"
const maxStringLength = 100

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

func prepareArray(n int) []keyValue {
	arr := make([]keyValue, 0, n)
	existing := make(map[string]bool)
	for i := 0; i < n; i++ {
		key := randomString(0)
		for existing[key] {
			key = randomString(0)
		}
		existing[key] = true
		arr = append(arr, keyValue{key, randomString(0)})
	}
	return arr
}

func BenchmarkSet(b *testing.B) {
	var n int
	for _, e := range []int{1, 2, 3} {
		n = 1
		for i := 0; i < e; i++ {
			n *= 10
		}
		b.Run(fmt.Sprintf("1e%d", e), func(b *testing.B) { benchmarkSet(b, n) })
	}
}

func benchmarkSet(b *testing.B, n int) {
	array := prepareArray(n)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d := structure.NewSipHashDict()
		for _, value := range array {
			d.Set(value.Key, value.Value)
		}
	}
	b.StopTimer()
}

func BenchmarkGet(b *testing.B) {
	var n int
	for _, e := range []int{1, 2, 3} {
		n = 1
		for i := 0; i < e; i++ {
			n *= 10
		}
		b.Run(fmt.Sprintf("1e%d", e), func(b *testing.B) { benchmarkGet(b, n) })
	}
}

func benchmarkGet(b *testing.B, n int) {
	array := prepareArray(n)
	d := structure.NewSipHashDict()
	for _, value := range array {
		d.Set(value.Key, value.Value)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, value := range array {
			val := d.Get(value.Key)
			if val != value.Value {
				b.Fatalf("Error getting element {%s, %v} from dictionary - got {%v}", value.Key, value.Value, val)
			}
		}
	}
	b.StopTimer()
}

func BenchmarkDelete(b *testing.B) {
	var n int
	for _, e := range []int{1, 2, 3} {
		n = 1
		for i := 0; i < e; i++ {
			n *= 10
		}
		b.Run(fmt.Sprintf("1e%d", e), func(b *testing.B) { benchmarkDelete(b, n) })
	}
}

func benchmarkDelete(b *testing.B, n int) {
	array := prepareArray(n)
	d := structure.NewSipHashDict()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		for _, value := range array {
			d.Set(value.Key, value.Value)
		}
		b.StartTimer()
		for _, value := range array {
			if d.Delete(value.Key) != nil {
				b.Fatalf("Error deleting element {%s} from dictionary", value.Key)
			}
		}
	}
}

func BenchmarkGoMapSet(b *testing.B) {
	var n int
	for _, e := range []int{1, 2, 3} {
		n = 1
		for i := 0; i < e; i++ {
			n *= 10
		}
		b.Run(fmt.Sprintf("1e%d", e), func(b *testing.B) { benchmarkGoMapSet(b, n) })
	}
}
func benchmarkGoMapSet(b *testing.B, n int) {
	array := prepareArray(n)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m := make(map[string]interface{})
		for _, value := range array {
			m[value.Key] = value.Value
		}
	}
	b.StopTimer()
}

func BenchmarkGoMapGet(b *testing.B) {
	var n int
	for _, e := range []int{1, 2, 3} {
		n = 1
		for i := 0; i < e; i++ {
			n *= 10
		}
		b.Run(fmt.Sprintf("1e%d", e), func(b *testing.B) { benchmarkGoMapGet(b, n) })
	}
}

func benchmarkGoMapGet(b *testing.B, n int) {
	array := prepareArray(n)
	m := make(map[string]interface{})
	for _, value := range array {
		m[value.Key] = value.Value
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, value := range array {
			if m[value.Key] != value.Value {
				b.Fatalf("Error getting element {%s, %v} from dictionary", value.Key, value.Value)
			}
		}
	}
	b.StopTimer()
}

func BenchmarkGoMapDelete(b *testing.B) {
	var n int
	for _, e := range []int{1, 2, 3} {
		n = 1
		for i := 0; i < e; i++ {
			n *= 10
		}
		b.Run(fmt.Sprintf("1e%d", e), func(b *testing.B) { benchmarkGoMapDelete(b, n) })
	}
}

func benchmarkGoMapDelete(b *testing.B, n int) {
	array := prepareArray(n)
	m := make(map[string]interface{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		for _, value := range array {
			m[value.Key] = value.Value
		}
		b.StartTimer()
		for _, value := range array {
			delete(m, value.Key)
			if m[value.Key] != nil {
				b.Fatalf("Error deleting element {%s} from dictionary", value.Key)
			}
		}
	}
}
