package bitmap

import (
	"sync/atomic"
	"unsafe"
)

// SetAtomicUint32 sets the target bit to the target value inside the uint32
// encded bitmap.
func SetAtomicUint32(bitmap []uint32, targetBit int, targetValue bool) {
	targetIndex := targetBit / 32
	BitOffset := targetBit % 32

	for {
		localValue := atomic.LoadUint32(&bitmap[targetIndex])
		targetBytes := (*[4]byte)(unsafe.Pointer(&localValue))[:]
		if Get(targetBytes, BitOffset) == targetValue {
			return
		}
		referenceValue := localValue
		Set(targetBytes, BitOffset, targetValue)
		if atomic.CompareAndSwapUint32(&bitmap[targetIndex], referenceValue, localValue) {
			break
		}
	}
}

// GetAtomicUint32 gets the target bit from an uint32 encoded bitmap.
func GetAtomicUint32(bitmap []uint32, targetBit int) bool {
	data := (*[4]byte)(unsafe.Pointer(&bitmap[targetBit/32]))[:]
	return Get(data, targetBit%32)
}
