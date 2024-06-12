package hashing

import (
	"encoding/binary"
	"sync"

	"github.com/dchest/siphash"
	"github.com/dmarro89/go-redis-hashtable/utility"
)

type IHasher interface {
	Digest(message string) uint64
}

type Sip24Hasher struct {
	key0 uint64
	key1 uint64
}

var customPool = sync.Pool{
	New: func() interface{} {
		buf := make([]byte, 0, 256)
		return &buf
	},
}

func NewSip24Hasher() IHasher {
	key0, key1 := Split(utility.GetRandomBytes())
	return &Sip24Hasher{
		key0: key0,
		key1: key1,
	}
}

func NewSip24BytesHasher(arr [16]byte) IHasher {
	key0, key1 := Split(arr)
	return &Sip24Hasher{
		key0: key0,
		key1: key1,
	}
}

func (h *Sip24Hasher) Digest(message string) uint64 {
	bufPtr := customPool.Get().(*[]byte)
	buf := *bufPtr

	if cap(buf) < len(message) {
		buf = make([]byte, len(message))
	} else {
		buf = buf[:len(message)]
	}

	copy(buf, message)

	hash := siphash.Hash(h.key0, h.key1, buf)

	customPool.Put(bufPtr)

	return hash
}

func Split(key [16]byte) (uint64, uint64) {
	key0 := binary.LittleEndian.Uint64(key[:8])
	key1 := binary.LittleEndian.Uint64(key[8:])
	return key0, key1
}
