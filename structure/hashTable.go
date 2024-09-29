package structure

type HashTable struct {
	table    []*DictEntry
	size     int64
	sizemask uint64
	used     int64
}

var errHashTableSize = "hashTable size cannot be negative"

// NewHashTable creates a new HashTable with the specified size.
//
// Parameters:
// - size: the size of the HashTable.
//
// Returns:
// - *HashTable: a pointer to the newly created HashTable.
func NewHashTable(size int64) *HashTable {
	if size < 0 {
		panic(errHashTableSize)
	}

	return &HashTable{
		table:    make([]*DictEntry, size),
		size:     size,
		sizemask: uint64(size - 1),
	}
}

// empty checks if the hash table is empty.
//
// Returns true if the hash table is empty, false otherwise.
func (ht *HashTable) empty() bool {
	return ht.size == 0
}
