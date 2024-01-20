package datastr

type HashTable struct {
	table    []*DictEntry
	size     int64
	sizemask int64
	used     int64
}

// NewHashTable creates a new HashTable with the specified size.
//
// Parameters:
// - size: the size of the HashTable.
//
// Returns:
// - *HashTable: a pointer to the newly created HashTable.
func NewHashTable(size int64) *HashTable {
	var sizemask int64
	table := []*DictEntry{}

	if size > 0 {
		table = make([]*DictEntry, size)
		sizemask = size - 1
	}

	return &HashTable{
		table:    table,
		size:     size,
		sizemask: sizemask,
	}
}

// empty checks if the hash table is empty.
//
// Returns true if the hash table is empty, false otherwise.
func (ht *HashTable) empty() bool {
	return ht.size == 0
}
