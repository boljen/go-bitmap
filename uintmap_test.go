package bitmap

import "testing"

func TestSetAtomicAlt(t *testing.T) {
	slice := make([]uint32, 100)
	if GetAtomicUint32(slice, 32) {
		t.Fatal("should return false")
	}
	SetAtomicUint32(slice, 32, true)

	if !GetAtomicUint32(slice, 32) {
		t.Fatal("should return true")
	}
}

func BenchmarkAtomicUint32(b *testing.B) {
	bm := make([]uint32, (1000*1000)/32)
	index := 0
	for i := 0; i < b.N; i++ {
		index += 50
		if index == 100000 {
			index = 0
		}
		SetAtomicUint32(bm, index, GetAtomicUint32(bm, index))

	}
}

func BenchmarkAtomicUint32_Parallel(b *testing.B) {
	bm := make([]uint32, (1000*1000)/32)

	b.RunParallel(func(pb *testing.PB) {
		index := 0
		for pb.Next() {
			index += 50
			if index == 100000 {
				index = 0
			}
			SetAtomicUint32(bm, index, GetAtomicUint32(bm, index))

		}
	})
}
