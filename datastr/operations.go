package datastr

import "fmt"

// Get returns the value associated with the given key in the dictionary.
//
// Parameters:
// - key: the key to look up in the dictionary.
//
// Return:
// - interface{}: the value associated with the key, or nil if the key is not found.
func (d *Dict) Get(key string) interface{} {
	d.rwmux.RLock()
	defer d.rwmux.RUnlock()
	entry := d.getEntry(key)
	if entry == nil {
		return nil
	}
	retval := entry.value
	return retval // copy avoids race conditions
}

// Set sets the value of a key in the dictionary.
//
// Parameters:
//   - key: the key to set the value for.
//   - value: the value to set.
//
// Returns:
//   - error: an error if the key already exists in the dictionary.
func (d *Dict) Set(key string, value interface{}) error {
	d.rwmux.Lock()
	defer d.rwmux.Unlock()
	entry := d.getEntry(key)
	if entry != nil {
		entry.value = value
		return nil
	}
	retval := d.add(key, value) // copy avoids race conditions
	return retval
}

// Delete deletes an entry from the dictionary.
//
// Parameters:
// - key: the key of the entry to be deleted.
//
// Returns:
// - error: if the entry is not found.
func (d *Dict) Delete(key string) error {
	d.rwmux.Lock()
	defer d.rwmux.Unlock()
	dictEntry := d.delete(key, false)
	if dictEntry == nil {
		return fmt.Errorf(`entry not found`)
	}
	//log.Printf("deleted key='%s'", key)
	return nil
}

