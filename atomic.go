package bitmap

import (
	"sync/atomic"
	"unsafe"
)

var oobPanic = "SetAtomic not allowed on a bitmapSlice of cap() < 4"

// SetAtomic is similar to Set except that it performs the operation atomically.
func SetAtomic(bitmap []byte, targetBit int, targetValue bool) {
	ov := (*[1]uint32)(unsafe.Pointer(&bitmap[targetBit/32]))[:]
	SetAtomicUint32(ov, targetBit%32, targetValue)
}

// SetAtomic is similar to Set except that it performs the operation atomically.
// It needs a bitmapSlice where the capacity is at least 4 bytes.
func _SetAtomic(bitmapSlice []byte, targetBit int, targetValue bool) {
	targetByteIndex := targetBit / 8
	targetBitIndex := targetBit % 8
	targetOffset := 0

	// SetAtomic needs to modify 4 bytes of data so we panic when the slice
	// doesn't have a capacity of at least 4 bytes.
	if cap(bitmapSlice) < 4 {
		panic(oobPanic)
	}

	// Calculate the Offset of the targetByte inside the 4-byte atomic batch.
	// This is needed to ensure that atomic operations can happen as long as
	// the bitmapSlice equals 4 bytes or more.
	if cap(bitmapSlice) < targetByteIndex+3 {
		targetOffset = cap(bitmapSlice) - targetByteIndex
	}

	// This gets a pointer to the memory of 4 bytes inside the bitmapSlice.
	// It stores this pointer as an *uint32 so that it can be used to
	// execute sync.atomic operations.
	targetBytePointer := (*uint32)(unsafe.Pointer(&bitmapSlice[targetByteIndex-targetOffset]))

	for {
		// localValue is a copy of the uint32 value at *targetBytePointer.
		// It's used to check whether the targetBit must be updated,
		// and if so, to construct the new value for targetBytePointer.
		localValue := atomic.LoadUint32(targetBytePointer)

		// This "neutralizes" the uint32 conversion by getting a pointer to the
		// 4-byte array stored undereneath the uint32.
		targetByteCopyPointer := (*[4]byte)(unsafe.Pointer(&localValue))

		// Work is done when targetBit is already set to targetValue.
		if GetBit(targetByteCopyPointer[targetOffset], targetBitIndex) == targetValue {
			return
		}

		// Modify the targetBit and update memory so that the targetBit is the only bit
		// that has been modified in the batch.
		referenceValue := localValue
		SetBitRef(&targetByteCopyPointer[targetOffset], targetBitIndex, targetValue)
		if atomic.CompareAndSwapUint32(targetBytePointer, referenceValue, localValue) {
			break
		}
	}
}
