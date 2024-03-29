package datastr

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDict(t *testing.T) {
	d := NewDict()
	assert.NotNil(t, d, "Failed to create a new dictionary")
	assert.Equal(t, 2, len(d.hashTables), "Missing two hashtables")
	assert.NotNil(t, d.mainTable(), "Failed to get the main table")
	assert.NotNil(t, d.rehashingTable(), "Failed to get the rehashing table")
	assert.Equal(t, -1, d.rehashidx, "Unexpected rehash index")
}

func TestMainTable(t *testing.T) {
	d := &Dict{}
	assert.Nil(t, d.mainTable(), "mainTable should be nil when hashTables is empty")
	d.hashTables[0] = NewHashTable(0)
	assert.NotNil(t, d.mainTable(), "mainTable should not be nil")
}

func TestRehashingTable(t *testing.T) {
	d := &Dict{}
	assert.Nil(t, d.rehashingTable(), "rehashingTable should be nil when hashTables is empty")
	d.hashTables[1] = NewHashTable(0)
	assert.NotNil(t, d.rehashingTable(), "rehashingTable should not be nil")
}

func TestKeyIndex(t *testing.T) {
	d := NewDict()
	d.mainTable().table = make([]*DictEntry, 4)
	d.mainTable().table[0] = NewDictEntry("mango", nil)
	d.mainTable().table[0].next = NewDictEntry("orange", nil)

	index := d.keyIndex("banana", []byte{1, 2, 3, 4})
	assert.Equal(t, 3, index, "Unexpected index for nonexistent key")

	fmt.Printf("%v", d.hashTables)
	index = d.keyIndex("orange", []byte{1, 2, 3, 4})
	assert.Equal(t, -1, index, "Unexpected index for nonexistent key")
}

func TestExpandIfNeeded(t *testing.T) {
	d := NewDict()

	// Test when rehashing is false and mainTable is nil
	d.expandIfNeeded()
	assert.NotNil(t, d.mainTable(), "mainTable should not be nil after first expansion")
	assert.Equal(t, len(d.mainTable().table), 4, "mainTable should have 4 entries with the first expansion")
	assert.Empty(t, d.rehashingTable().table, "rehashingTable should be nil after first expansion")

	// Test when rehashing is false and mainTable is not nil
	d.expandIfNeeded()
	assert.Equal(t, len(d.mainTable().table), 4, "mainTable should have 4 entries with the first expansion")
	assert.True(t, d.mainTable().used < d.mainTable().size, "Unexpected expansion when not needed")

	// Trigger rehashing by setting used equal to size
	d.mainTable().used = d.mainTable().size
	d.expandIfNeeded()
	assert.NotEmpty(t, d.rehashingTable().table, "rehashingTable should not be nil during rehashing")
	assert.Equal(t, len(d.rehashingTable().table), 8, "rehashingTable should have 8 entries with the first expansion")
	assert.Equal(t, len(d.mainTable().table), 4, "mainTable should have 4 entries with the first expansion")
}

func TestExpand(t *testing.T) {
	d := NewDict()

	// Test when rehashing is true
	d.rehashidx = 0
	d.expand(10)
	assert.Empty(t, d.hashTables[1].table, "Expansion should be skipped during rehashing")
	assert.Equal(t, d.rehashidx, 0)

	// Test when rehashing is false and used is greater than newSize
	d.rehashidx = -1
	d.mainTable().used = 10
	d.expand(5)
	assert.Empty(t, d.hashTables[1].table, "Expansion should be skipped when used is greater than newSize")
	assert.Equal(t, d.rehashidx, -1)

	// Test a valid expansion for the main table
	assert.Empty(t, d.hashTables[0].table, "Main Table should be nil")
	d.mainTable().used = 5
	d.expand(10)
	assert.NotNil(t, d.hashTables[0].table, "Unexpected expansion")
	assert.Empty(t, d.hashTables[1].table, "Unexpected expansion")
	assert.Equal(t, len(d.mainTable().table), 16, "mainTable should have 16 entries with the first expansion")
	assert.Equal(t, -1, d.rehashidx, "Unexpected rehash index after expansion")

	// Test a valid expansion for the rehashing table
	d.mainTable().used = 18
	d.expand(24)
	assert.NotNil(t, d.hashTables[0].table, "Unexpected expansion for main table")
	assert.NotNil(t, d.hashTables[1].table, "Missing expansio for rehashing table")
	assert.Equal(t, len(d.mainTable().table), 16, "mainTable should have 16 entries with the first expansion")
	assert.Equal(t, len(d.rehashingTable().table), 32, "rehashingTable should have 32 entries with the first expansion")
	assert.Equal(t, 0, d.rehashidx, "Unexpected rehash index after expansion")
}

func TestNextPower(t *testing.T) {
	size := nextPower(6)
	assert.Equal(t, int64(8), size, "Unexpected next power of 2")

	size = nextPower(16)
	assert.Equal(t, int64(16), size, "Unexpected next power of 2 for an already power of 2 input")
}

func TestSipHashDigest(t *testing.T) {
	randomBytes := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	hash := sipHashDigest(randomBytes, "test_key_1_2_3_4")
	assert.NotZero(t, hash, "Unexpected SIP hash digest")
}

func TestAdd(t *testing.T) {
	dictionary := NewDict()

	key1 := "keyTest"
	value1 := 123
	err := dictionary.add(key1, value1)
	assert.NoError(t, err, "Unexpected error adding a new key")

	entry := dictionary.getEntry(key1)
	assert.NotNil(t, entry, "Key %s was not added correctly", key1)
	assert.Equal(t, value1, entry.value, "Value associated with key %s is incorrect", key1)

	err = dictionary.add(key1, "newValue")
	assert.Error(t, err, "There should be an error for adding an existing key")
	assert.EqualError(t, err, fmt.Sprintf(`unexpectedly found an entry with the same key when trying to add #{ %s } / #{ %s }`, key1, "newValue"), "Unexpected error")

	entry = dictionary.getEntry(key1)
	assert.NotNil(t, entry, "Key %s was not added correctly", key1)
	assert.Equal(t, value1, entry.value, "Value associated with key %s is incorrect", key1)

	// Adding more entries for rehashing test
	for i := 0; i < 4; i++ {
		key := fmt.Sprintf("key%d", i)
		value := rand.Intn(100)
		err := dictionary.add(key, value)
		assert.NoError(t, err, "Unexpected error adding a new key")

		entry := dictionary.getEntry(key)
		assert.NotNil(t, entry, "Key %s was not added correctly", key)
		assert.Equal(t, value, entry.value, "Value associated with key %s is incorrect", key)
	}
}

func TestRehash(t *testing.T) {
	d := NewDict()

	//Not rehashing
	result := d.rehash(1)
	assert.Zero(t, result, "Unexpected result when not rehashing")

	// Test rehashing an empty main table
	d.rehashidx = 0
	result = d.rehash(1)
	assert.Zero(t, result, "Unexpected result when rehashing an empty table")
	assert.Equal(t, d.rehashidx, -1)

	//Rehashing last element from rehashing table
	d = NewDict()
	d.rehashingTable().table = make([]*DictEntry, 1)
	d.rehashingTable().table[0] = NewDictEntry("key-test", "value-test")
	d.add("key1", "value1")
	d.rehashidx = 0
	result = d.rehash(1)
	assert.Zero(t, result, "Unexpected result when rehashing an empty table")
	assert.Equal(t, d.rehashidx, -1)
	assert.Equal(t, d.mainTable().table[0].key, "key1")
	assert.Equal(t, d.mainTable().table[0].value, "value1")
	assert.Equal(t, d.mainTable().table[0].next.key, "key-test")
	assert.Equal(t, d.mainTable().table[0].next.value, "value-test")
	assert.Empty(t, d.rehashingTable().table)
}

func TestRehashing(t *testing.T) {
	d := NewDict()

	assert.False(t, d.rehashing(), "Unexpected rehashing status when rehashing is false")
	d.rehashidx = 0
	assert.True(t, d.rehashing(), "Unexpected rehashing status when rehashing is true")
}

func TestGetEntry(t *testing.T) {
	// Test getEntry method
	d := NewDict()

	// Test when both tables are empty
	entry := d.getEntry("nonexistent_key")
	assert.Nil(t, entry, "Unexpected entry for nonexistent key when tables are empty")

	// Test when rehashing and mainTable is empty
	d.rehashidx = 0
	entry = d.getEntry("nonexistent_key")
	assert.Nil(t, entry, "Unexpected entry for nonexistent key when rehashing and mainTable is empty")

	// Test when key is in rehashingTable
	d.rehashidx = -1
	d.add("key1", "value1")
	entry = d.getEntry("key1")
	assert.NotNil(t, entry, "Expected entry for key1 in mainTable")
	assert.Equal(t, "key1", entry.key, "Unexpected key in entry")
	assert.Equal(t, "value1", entry.value, "Unexpected key in entry")
}

func TestDelete(t *testing.T) {
	// Test delete method
	d := NewDict()

	// Test deleting a key that does not exist
	deletedEntry := d.delete("nonexistent_key")
	assert.Nil(t, deletedEntry, "Unexpected entry deleted for nonexistent key")

	// Test deleting a key from an empty table
	deletedEntry = d.delete("key1")
	assert.Nil(t, deletedEntry, "Unexpected entry deleted from an empty table")

	// Test deleting a key from mainTable
	d.add("key1", "value1")
	deletedEntry = d.delete("key1")
	assert.NotNil(t, deletedEntry, "Expected entry deleted from mainTable")
	assert.Equal(t, "key1", deletedEntry.key, "Unexpected key in deleted entry")
}
