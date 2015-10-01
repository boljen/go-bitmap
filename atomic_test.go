package bitmap

import "testing"

func TestSetAtomic(t *testing.T) {
	slice := make([]byte, 4)
	SetAtomic(slice, 31, true)
	if !Get(slice, 31) {
		t.Fatal("should return true")
	}
}

func BenchmarkAtomic(b *testing.B) {
	bm := NewConcurrent(10000)
	index := 0
	for i := 0; i < b.N; i++ {
		index++
		if index == 1000 {
			index = 0
		}
		bm.Set(index, !bm.Get(index))
	}
}

func BenchmarkAtomic_Parallel(b *testing.B) {
	bm := NewConcurrent(1000 * 1000)

	b.RunParallel(func(pb *testing.PB) {
		index := 0
		for pb.Next() {
			index += 50
			if index == 100000 {
				index = 0
			}
			bm.Set(index, !bm.Get(index))

		}
	})
}
