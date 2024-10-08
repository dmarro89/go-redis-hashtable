package structure

import (
	"fmt"

	"github.com/dmarro89/go-redis-hashtable/hashing"
)

const (
	INITIAL_SIZE = int64(4)
	MAX_SIZE     = 1 << 63
)

type IDict interface {
	Set(key string, value string) error
	Get(key string) string
	Delete(key string) error
	GetAllItems() map[string]string
}

type Dict struct {
	hashTables [2]*HashTable
	rehashidx  int
	hasher     hashing.IHasher
}

// NewSipHashDict returns a new instance of Dict.
//
// The function does not take any parameters.
// It returns a pointer to Dict.
func NewSipHashDict() IDict {
	return &Dict{
		hashTables: [2]*HashTable{NewHashTable(0), NewHashTable(0)},
		rehashidx:  -1,
		hasher:     hashing.NewSip24Hasher(),
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
	if d.isRehashing() || d.mainTable().used > newSize {
		return
	}

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

// keyIndex returns the index of the given key in the dictionary.
//
// It takes in a key string and randomBytes []byte as parameters.
// It returns an integer representing the index of the key in the dictionary.
func (d *Dict) keyIndex(key string) int {
	d.expandIfNeeded()
	hash := d.hasher.Digest(key)

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
func (d *Dict) add(key string, value string) error {
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
			idx := d.hasher.Digest(entry.key) & d.rehashingTable().sizemask

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

	hash := d.hasher.Digest(key)

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
func (d *Dict) delete(key string) *DictEntry {
	if d.mainTable().used == 0 && d.rehashingTable().used == 0 {
		return nil
	}

	if d.isRehashing() {
		d.rehashStep()
	}

	hash := d.hasher.Digest(key)

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

// Get returns the value associated with the given key in the dictionary.
//
// Parameters:
// - key: the key to look up in the dictionary.
//
// Return:
// - interface{}: the value associated with the key, or nil if the key is not found.
func (d *Dict) Get(key string) string {
	entry := d.getEntry(key)
	if entry == nil {
		return ""
	}
	return entry.value
}

// Set sets the value of a key in the dictionary.
//
// Parameters:
//   - key: the key to set the value for.
//   - value: the value to set.
//
// Returns:
//   - error: an error if the key already exists in the dictionary.
func (d *Dict) Set(key string, value string) error {
	entry := d.getEntry(key)
	if entry != nil {
		entry.value = value
		return nil
	}
	return d.add(key, value)
}

// Delete deletes an entry from the dictionary.
//
// Parameters:
// - key: the key of the entry to be deleted.
//
// Returns:
// - error: if the entry is not found.
func (d *Dict) Delete(key string) error {
	dictEntry := d.delete(key)
	if dictEntry == nil {
		return fmt.Errorf(`entry not found`)
	}
	return nil
}

// GetAllKeys retrieves all keys from the hash table.
// It iterates over both hash tables in the Dict struct (main table and rehashing table).
// Each bucket may contain a linked list of entries (DictEntry) due to hash collisions,
// so it traverses through these linked lists to collect all keys.
// This function supports the dynamic resizing and rehashing mechanism.
func (d *Dict) GetAllItems() map[string]string {
	items := make(map[string]string)

	// Iterate over both hash tables (HashTable[2])
	for _, hashtable := range d.hashTables {
		if hashtable != nil {
			// Iterate through each bucket in the hash table
			for _, entry := range hashtable.table {
				// Traverse the linked list at each index to get all keys
				for entry != nil {
					items[entry.key] = entry.value
					entry = entry.next
				}
			}
		}
	}

	return items
}
