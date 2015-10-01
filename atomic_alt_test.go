package bitmap

import "testing"

import (
	"sync/atomic"
	"unsafe"
)

// SetAtomicAlt is similar to Set except that it performs the operation atomically.
// It needs a []uint32 slice that stores the bitmap.
func SetAtomicAlt(bitmapSlice []uint32, targetBit int, targetValue bool) {
	targetIndex := targetBit / 32
	BitOffset := targetBit % 32

	for {
		localValue := atomic.LoadUint32(&bitmapSlice[targetIndex])

		targetBytes := (*[4]byte)(unsafe.Pointer(&localValue))[:]

		// Work is done when targetBit is already set to targetValue.
		if Get(targetBytes, BitOffset) == targetValue {
			return
		}

		// Modify the targetBit and update memory so that the targetBit is the only bit
		// that has been modified in the batch.
		referenceValue := localValue
		Set(targetBytes, BitOffset, targetValue)
		if atomic.CompareAndSwapUint32(&bitmapSlice[targetIndex], referenceValue, localValue) {
			break
		}
	}
}

func GetAtomicAlt(slice []uint32, targetBit int) bool {
	data := (*[4]byte)(unsafe.Pointer(&slice[targetBit/32]))[:]
	return Get(data, targetBit%32)
}

func TestSetAtomicAlt(t *testing.T) {
	slice := make([]uint32, 100)
	SetAtomicAlt(slice, 32, true)
	data := (*[4]byte)(unsafe.Pointer(&slice[1]))[:]

	if !Get(data, 0) {
		t.Fatal("should return true")
	}
}

func BenchmarkAlt(b *testing.B) {
	bm := make([]uint32, (1000*1000)/32)

	b.RunParallel(func(pb *testing.PB) {
		index := 0
		for pb.Next() {
			index += 50
			if index == 100000 {
				index = 0
			}
			SetAtomicAlt(bm, index, GetAtomicAlt(bm, index))

		}
	})
}
