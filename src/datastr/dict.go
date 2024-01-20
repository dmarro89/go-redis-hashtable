package datastr

import (
	"fmt"
	"go-redis-hashtable/src/utilities"

	"github.com/dchest/siphash"
)

const (
	INITIAL_SIZE = int64(4)
	MAX_SIZE     = 1 << 63
)

type Dict struct {
	hashTables [2]*HashTable
	rehashidx  int
}

// NewDict returns a new instance of Dict.
//
// The function does not take any parameters.
// It returns a pointer to Dict.
func NewDict() *Dict {
	return &Dict{
		hashTables: [2]*HashTable{NewHashTable(0), NewHashTable(0)},
		rehashidx:  -1,
	}
}

// mainTable returns the main hash table of the Dict.
//
// No parameters.
// Returns a pointer to a HashTable.
func (d *Dict) mainTable() *HashTable {
	return d.hashTables[0]
}

// rehashingTable returns the HashTable at index 1 of the Dict.
//
// No parameters.
// Returns *HashTable.
func (d *Dict) rehashingTable() *HashTable {
	return d.hashTables[1]
}

// nextPower calculates the next power of 2 greater than the given size.
//
// Parameters:
// - size: the size for which we want to find the next power of 2.
//
// Return type:
// - int64: the next power of 2 greater than the given size.
func nextPower(size int64) int64 {
	i := INITIAL_SIZE
	for ; i < size; i *= 2 {
	}
	return i
}

// expand expands the dictionary to a new size if necessary.
//
// newSize: the new size to expand the dictionary to.
// The function does not return anything.
func (d *Dict) expand(newSize int64) {
	if d.rehashing() || d.mainTable().used > newSize {
		return
	}

	realSize := nextPower(newSize)

	if realSize != d.mainTable().size {
		newHashTable := NewHashTable(realSize)

		if d.mainTable() == nil || len(d.mainTable().table) == 0 {
			d.hashTables[0] = newHashTable
		} else {
			d.hashTables[1] = newHashTable
			d.rehashidx = 0
		}
	}
}

// expandIfNeeded checks if the dictionary needs to be expanded and performs the expansion if necessary.
//
// No parameters.
// No return values.
func (d *Dict) expandIfNeeded() {
	if d.rehashing() {
		return
	}

	if d.mainTable() == nil || len(d.mainTable().table) == 0 {
		d.expand(INITIAL_SIZE)
	} else if d.mainTable().used >= d.mainTable().size {
		newSize := int64(d.mainTable().used * 2)
		d.expand(newSize)
	}
}

// sipHashDigest calculates the SipHash-2-4 digest of the given random bytes using the provided key.
//
// randomBytes: The random bytes to calculate the digest for.
// key: The key used for the SipHash-2-4 algorithm.
// Returns the calculated 64-bit SipHash-2-4 digest.
func sipHashDigest(randomBytes []byte, key string) uint64 {
	key16 := make([]byte, 16)
	copy(key16, key)
	keyBytes := []byte(key16)
	h := siphash.New(keyBytes)
	h.Write(randomBytes)
	return h.Sum64()
}

// keyIndex returns the index of the given key in the dictionary.
//
// It takes in a key string and randomBytes []byte as parameters.
// It returns an integer representing the index of the key in the dictionary.
func (d *Dict) keyIndex(key string, randomBytes []byte) int {
	d.expandIfNeeded()

	hash := sipHashDigest(randomBytes, key)

	var index int
	for i := 0; i <= 1; i++ {
		hashTable := d.hashTables[i]
		index = int(hash & uint64(hashTable.sizemask))

		for entry := hashTable.table[index]; entry != nil; entry = entry.next {
			if entry.key == key {
				return -1
			}
		}

		if index == -1 || !d.rehashing() {
			break
		}
	}

	return index
}

// add adds a key-value pair to the dictionary.
//
// Parameters:
// - key: The key to add.
// - value: The value associated with the key.
//
// Returns:
// - error: An error if the key already exists in the dictionary.
func (d *Dict) add(key string, value interface{}) error {
	index := d.keyIndex(key, utilities.GetRandomBytes())

	if index == -1 {
		return fmt.Errorf(`unexpectedly found an entry with the same key when trying to add #{ %s } / #{ %s }`, key, value)
	}

	hashTable := d.mainTable()
	if d.rehashing() {
		hashTable = d.rehashingTable()
	}

	entry := hashTable.table[index]

	for entry != nil && entry.key != key {
		entry = entry.next
	}

	if entry == nil {
		entry = NewDictEntry(key, value)
		entry.next = hashTable.table[index]
		hashTable.table[index] = entry
		hashTable.used++
	}

	return nil
}

// rehashStep returns the result of calling the rehash function on the Dict object with an argument of 1.
//
// No parameters.
// Returns an integer.
func (d *Dict) rehashStep() int {
	return d.rehash(1)
}

// rehash rehashes the dictionary with a new size.
//
// n is the new size of the dictionary.
// Returns 0 if the rehashing is not in progress.
// Returns 1 if the rehashing is in progress.
func (d *Dict) rehash(n int) int {
	emptyVisits := n * 10
	if !d.rehashing() {
		return 0
	}

	for n > 0 && d.mainTable().used != 0 {
		n--

		var entry *DictEntry

		for len(d.mainTable().table) == 0 || d.mainTable().table[d.rehashidx] == nil {
			d.rehashidx++
			emptyVisits--
			if emptyVisits == 0 {
				return 1
			}
		}

		entry = d.mainTable().table[d.rehashidx]

		for entry != nil {
			nextEntry := entry.next
			randomBytes := utilities.GetRandomBytes()
			idx := sipHashDigest(randomBytes, entry.key) & uint64(d.rehashingTable().sizemask)

			entry.next = d.rehashingTable().table[idx]
			d.rehashingTable().table[idx] = entry
			d.mainTable().used--
			d.rehashingTable().used++
			entry = nextEntry
		}

		d.mainTable().table[d.rehashidx] = nil
		d.rehashidx++
	}

	if d.mainTable().used == 0 {
		d.hashTables[0] = d.rehashingTable()
		d.hashTables[1] = NewHashTable(0)
		d.rehashidx = -1
		return 0
	}

	return 1
}

// rehashing checks if the rehash index of the Dict struct is not equal to -1.
//
// It does not take any parameters.
// It returns a boolean value indicating whether the rehash index is not equal to -1.
func (d *Dict) rehashing() bool {
	return d.rehashidx != -1
}

// getEntry returns the DictEntry associated with the given key in the Dict.
//
// Parameters:
// - key: the key to search for in the Dict.
//
// Return:
// - *DictEntry: the DictEntry associated with the given key, or nil if not found.
func (d *Dict) getEntry(key string) *DictEntry {
	if d.mainTable().used == 0 && d.rehashingTable().used == 0 {
		return nil
	}

	if d.rehashing() {
		d.rehashStep()
	}

	randomBytes := utilities.GetRandomBytes()
	hash := sipHashDigest(randomBytes, key)

	for ind, hashTable := range []*HashTable{d.mainTable(), d.rehashingTable()} {
		if hashTable == nil || len(hashTable.table) == 0 || (ind == 1 && !d.rehashing()) {
			continue
		}

		index := hash & uint64(hashTable.sizemask)
		entry := hashTable.table[index]

		for entry != nil {
			if entry.key == key {
				return entry
			}
			entry = entry.next
		}
	}

	return nil
}

// delete deletes a key from the dictionary and returns the corresponding value.
//
// Parameters:
// - key: the key to be deleted from the dictionary.
//
// Return:
// - *DictEntry: the deleted DictEntry if found, otherwise nil.
func (d *Dict) delete(key string) *DictEntry {
	if d.mainTable().used == 0 && d.rehashingTable().used == 0 {
		return nil
	}

	if d.rehashing() {
		d.rehashStep()
	}

	randomBytes := utilities.GetRandomBytes()
	hash := sipHashDigest(randomBytes, key)

	for i, hashTable := range []*HashTable{d.mainTable(), d.rehashingTable()} {
		if hashTable == nil || (i == 1 && !d.rehashing()) {
			continue
		}
		index := hash & uint64(hashTable.sizemask)
		entry := hashTable.table[index]
		var previousEntry *DictEntry

		for entry != nil {
			if entry.key == key {
				if previousEntry != nil {
					previousEntry.next = entry.next
				} else {
					hashTable.table[index] = entry.next
				}
				hashTable.used--
				return entry
			}
			previousEntry = entry
			entry = entry.next
		}
	}

	return nil
}
