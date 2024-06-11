package datastr

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"github.com/dchest/siphash"
	"sync"
	"time"
)

const (
	INITIAL_SIZE = int64(4)
	MAX_SIZE     = 1 << 63
)

type Dict struct {
	hashTables [2]*HashTable
	rehashidx  int
	rwmux sync.RWMutex
	randomBytes [16]byte
	once sync.Once
	key0 uint64
	key1 uint64
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
	if size <= INITIAL_SIZE {
		return INITIAL_SIZE
	}

	size--
	size |= size >> 1
	size |= size >> 2
	size |= size >> 4
	size |= size >> 8
	size |= size >> 16
	size |= size >> 32

	return size + 1
}

// expand expands the dictionary to a new size if necessary.
//
// newSize: the new size to expand the dictionary to.
// The function does not return anything.
func (d *Dict) expand(newSize int64) {
	isrehashing := d.isRehashing()
	istablefull := d.mainTable().used > newSize
	if isrehashing || istablefull {
		//log.Printf("dict.expand return1 newSize=%d isrehashing=%t istablefull=%t", newSize, isrehashing, istablefull)
		return
	}

	//log.Printf("dict.expand newSize=%d isrehashing=%t istablefull=%t", newSize, isrehashing, istablefull)

	nextSize := nextPower(newSize)
	if d.mainTable().used >= nextSize {
		return
	}

	newHashTable := NewHashTable(nextSize)

	if d.mainTable() == nil || len(d.mainTable().table) == 0 {
		*d.mainTable() = *newHashTable
		return
	}

	*d.rehashingTable() = *newHashTable
	d.rehashidx = 0
}

// expandIfNeeded checks if the dictionary needs to be expanded and performs the expansion if necessary.
//
// No parameters.
// No return values.
func (d *Dict) expandIfNeeded() {
	if d.isRehashing() {
		return
	}

	if d.mainTable() == nil || len(d.mainTable().table) == 0 {
		d.expand(INITIAL_SIZE)
	} else if d.mainTable().used >= d.mainTable().size {
		newSize := int64(d.mainTable().used * 2)
		d.expand(newSize)
	}
}

func (d *Dict) split(key [16]byte) (uint64, uint64) {
	if len(key) == 0 || len(key) < 16 {
		//d.logger.Error("ERROR split len(key)=%d", len(key))
		return 0, 0
	}
	key0 := binary.LittleEndian.Uint64(key[:8])
	key1 := binary.LittleEndian.Uint64(key[8:16])
	return key0, key1
}

// sipHashDigest calculates the SipHash-2-4 digest of the given message using the provided key.
//
// Parameters:
// - key: The key used for the SipHash-2-4 algorithm. It should be a byte slice of length 16.
// - message: The message to calculate the digest for.
//
// Returns:
// - uint64: The calculated 64-bit SipHash-2-4 digest.
func (d *Dict) sipHashDigest(message string) uint64 {
	//log.Printf("sipHashDigest msg='%s' key0=%d key1=%d", message, d.key0, d.key1)
	return siphash.Hash(d.key0, d.key1, []byte(message))
}

// keyIndex returns the index of the given key in the dictionary.
//
// It takes in a key string and randomBytes []byte as parameters.
// It returns an integer representing the index of the key in the dictionary.
func (d *Dict) keyIndex(key string) int {
	//log.Printf("keyIndex(key=%d='%s'", len(key), key)
	d.expandIfNeeded()
	hash := d.sipHashDigest(key)

	var index int
	for i := 0; i <= 1; i++ {
		hashTable := d.hashTables[i]
		index = int(hash & hashTable.sizemask)

		for entry := hashTable.table[index]; entry != nil; entry = entry.next {
			if entry.key == key {
				return -1
			}
		}

		if index == -1 || !d.isRehashing() {
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
	//log.Printf("add(key=%d='%s' value='%#v'", len(key), key, value)
	index := d.keyIndex(key)

	if index == -1 {
		return fmt.Errorf(`unexpectedly found an entry with the same key when trying to add #{ %s } / #{ %s }`, key, value)
	}

	hashTable := d.mainTable()
	if d.isRehashing() {
		d.rehashStep()
		hashTable = d.mainTable()
		if d.isRehashing() {
			hashTable = d.rehashingTable()
		}
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
func (d *Dict) rehashStep() {
	d.rehash(1)
}

// rehash rehashes the dictionary with a new size.
//
// n is the new size of the dictionary.
// Returns 0 if the rehashing is not in progress.
// Returns 1 if the rehashing is in progress.
func (d *Dict) rehash(n int) {
	emptyVisits := n * 10
	if !d.isRehashing() {
		return
	}

	for n > 0 && d.mainTable().used != 0 {
		n--

		var entry *DictEntry

		for len(d.mainTable().table) == 0 || d.mainTable().table[d.rehashidx] == nil {
			d.rehashidx++
			emptyVisits--
			if emptyVisits == 0 {
				return
			}
		}

		entry = d.mainTable().table[d.rehashidx]

		for entry != nil {
			nextEntry := entry.next
			idx := d.sipHashDigest(entry.key) & d.rehashingTable().sizemask

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
		return
	}
}

// isRehashing checks if the rehash index of the Dict struct is not equal to -1.
//
// It does not take any parameters.
// It returns a boolean value indicating whether the rehash index is not equal to -1.
func (d *Dict) isRehashing() bool {
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

	hash := d.sipHashDigest(key)

	for ind, hashTable := range []*HashTable{d.mainTable(), d.rehashingTable()} {
		if hashTable == nil || len(hashTable.table) == 0 || (ind == 1 && !d.isRehashing()) {
			continue
		}

		index := hash & hashTable.sizemask
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
func (d *Dict) delete(key string, dolock bool) *DictEntry {
	if dolock {
		d.rwmux.Lock()
		defer d.rwmux.Unlock()
	}
	if d.mainTable().used == 0 && d.rehashingTable().used == 0 {
		return nil
	}

	if d.isRehashing() {
		d.rehashStep()
	}

	hash := d.sipHashDigest(key)

	for i, hashTable := range []*HashTable{d.mainTable(), d.rehashingTable()} {
		if hashTable == nil || (i == 1 && !d.isRehashing()) {
			continue
		}
		index := hash & hashTable.sizemask
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

// GenerateRandomBytes generates a fixed slice of 16 random bytes
//
// Parameters:
// - none
//
// Returns:
// - none
func (d *Dict) GenerateRandomBytes() {
	d.once.Do(func() {
		rand.Seed(time.Now().UnixNano())
		cs := "0123456789abcdef"
		for i := 0; i < 16; i++ {
			d.randomBytes[i] = cs[rand.Intn(len(cs))]
		}
		d.key0, d.key1 = d.split(d.randomBytes)
	})
}



