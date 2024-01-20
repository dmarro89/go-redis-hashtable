package datastr

type DictEntry struct {
	next  *DictEntry
	key   string
	value interface{}
}

// NewDictEntry creates a new DictEntry with the given key and value.
//
// Parameters:
// - key: a string representing the key of the entry.
// - value: an interface{} representing the value of the entry.
//
// Returns:
// - *DictEntry: a pointer to the newly created DictEntry.
func NewDictEntry(key string, value interface{}) *DictEntry {
	return &DictEntry{
		key:   key,
		value: value,
		next:  nil,
	}
}
