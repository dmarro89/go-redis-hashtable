package structure

type DictEntry struct {
	next  *DictEntry
	key   string
	value string
}

// NewDictEntry creates a new DictEntry with the given key and value.
//
// Parameters:
// - key: a string representing the key of the entry.
// - value: an interface{} representing the value of the entry.
//
// Returns:
// - *DictEntry: a pointer to the newly created DictEntry.
func NewDictEntry(key string, value string) *DictEntry {
	return &DictEntry{
		key:   key,
		value: value,
		next:  nil,
	}
}
