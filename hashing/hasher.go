package hashing

import (
	"encoding/binary"

	"github.com/dchest/siphash"
	"github.com/dmarro89/go-redis-hashtable/utility"
)

type IHasher interface {
	Digest(message string) uint64
}

type Sip24Hasher struct {
	Key0 uint64
	Key1 uint64
}

var key0, key1 = Split(utility.GetRandomBytes())
var globalSip24Hasher = &Sip24Hasher{Key0: key0, Key1: key1}

func NewSip24Hasher() IHasher {
	return globalSip24Hasher
}

func (h *Sip24Hasher) Digest(message string) uint64 {
	return siphash.Hash(h.Key0, h.Key1, []byte(message))
}

func Split(key [16]byte) (uint64, uint64) {
	key0 := binary.LittleEndian.Uint64(key[:8])
	key1 := binary.LittleEndian.Uint64(key[8:])
	return key0, key1
}
