// Package bitmap implements (thread-safe) bitmap functions and abstractions.
//
// Installation
//
// 	  go get github.com/boljen/go-bitmap
package bitmap

import "sync"

var (
	tA = [8]byte{1, 2, 4, 8, 16, 32, 64, 128}
	tB = [8]byte{254, 253, 251, 247, 239, 223, 191, 127}
)

func dataOrCopy(d []byte, c bool) []byte {
	if !c {
		return d
	}
	ndata := make([]byte, len(d))
	copy(ndata, d)
	return ndata
}

// NewSlice creates a new byteslice with length l (in bits).
// The actual size in bits might be up to 7 bits larger because
// they are stored in a byteslice.
func NewSlice(l int) []byte {
	remainder := l % 8
	if remainder != 0 {
		remainder = 1
	}
	return make([]byte, l/8+remainder)
}

// Get returns the value of bit i from map m.
// It doesn't check the bounds of the slice.
func Get(m []byte, i int) bool {
	return m[i/8]&tA[i%8] != 0
}

// Set sets bit i of map m to value v.
// It doesn't check the bounds of the slice.
func Set(m []byte, i int, v bool) {
	index := i / 8
	bit := i % 8
	if v {
		m[index] = m[index] | tA[bit]
	} else {
		m[index] = m[index] & tB[bit]
	}
}

// GetBit returns the value of bit i of byte b.
// The bit index must be between 0 and 7.
func GetBit(b byte, i int) bool {
	return b&tA[i] != 0
}

// SetBit sets bit i of byte b to value v.
// The bit index must be between 0 and 7.
func SetBit(b byte, i int, v bool) byte {
	if v {
		return b | tA[i]
	}
	return b & tB[i]
}

// SetBitRef sets bit i of byte *b to value v.
func SetBitRef(b *byte, i int, v bool) {
	if v {
		*b = *b | tA[i]
	} else {
		*b = *b & tB[i]
	}
}

// Len returns the length (in bits) of the provided byteslice.
// It will always be a multipile of 8 bits.
func Len(m []byte) int {
	return len(m) * 8
}

// Bitmap is a byteslice with bitmap functions.
// Creating one form existing data is as simple as bitmap := Bitmap(data).
type Bitmap []byte

// New creates a new Bitmap instance with length l (in bits).
func New(l int) Bitmap {
	return NewSlice(l)
}

// Len wraps around the Len function.
func (b Bitmap) Len() int {
	return Len(b)
}

// Get wraps around the Get function.
func (b Bitmap) Get(i int) bool {
	return Get(b, i)
}

// Set wraps around the Set function.
func (b Bitmap) Set(i int, v bool) {
	Set(b, i, v)
}

// Data returns the data of the bitmap.
// If copy is false the actual underlying slice will be returned.
func (b Bitmap) Data(copy bool) []byte {
	return dataOrCopy(b, copy)
}

// Threadsafe implements thread-safe read- and write locking for the bitmap.
type Threadsafe struct {
	bm Bitmap
	mu sync.RWMutex
}

// TSFromData creates a new Threadsafe using the provided data.
// If copy is true the actual slice will be used.
func TSFromData(data []byte, copy bool) *Threadsafe {
	return &Threadsafe{
		bm: Bitmap(dataOrCopy(data, copy)),
	}
}

// NewTS creates a new Threadsafe instance.
func NewTS(length int) *Threadsafe {
	return &Threadsafe{
		bm: New(length),
	}
}

// Data returns the data of the bitmap.
// If copy is false the actual underlying slice will be returned.
func (b *Threadsafe) Data(copy bool) []byte {
	b.mu.RLock()
	data := dataOrCopy(b.bm, copy)
	b.mu.RUnlock()
	return data
}

// Len wraps around the Len function.
func (b Threadsafe) Len() int {
	b.mu.RLock()
	l := b.bm.Len()
	b.mu.RUnlock()
	return l
}

// Get wraps around the Get function.
func (b Threadsafe) Get(i int) bool {
	b.mu.RLock()
	v := b.bm.Get(i)
	b.mu.RUnlock()
	return v
}

// Set wraps around the Set function.
func (b Threadsafe) Set(i int, v bool) {
	b.mu.Lock()
	b.bm.Set(i, v)
	b.mu.Unlock()
}

// Concurrent is a bitmap implementation that achieves thread-safety
// using atomic operations along with some unsafe.
// It performs atomic operations on 32bits of data.
type Concurrent []byte

// NewConcurrent returns a concurrent bitmap.
// It will create a bitmap
func NewConcurrent(l int) Concurrent {
	remainder := l % 8
	if remainder != 0 {
		remainder = 1
	}
	return make([]byte, l/8+remainder, l/8+remainder+3)
}

// Get wraps around the Get function.
func (c Concurrent) Get(b int) bool {
	return Get(c, b)
}

// Set wraps around the SetAtomic function.
func (c Concurrent) Set(b int, v bool) {
	SetAtomic(c, b, v)
}

// Len wraps around the Len function.
func (c Concurrent) Len() int {
	return Len(c)
}

// Data returns the data of the bitmap.
// If copy is false the actual underlying slice will be returned.
func (c Concurrent) Data(copy bool) []byte {
	return dataOrCopy(c, copy)
}
