package bitmap

import "testing"

func TestSetAtomicPanic(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil || err != oobPanic {
			t.Fatal("expected a recover panic")
		}
	}()

	slice := make([]byte, 3)
	SetAtomic(slice, 0, true)
	t.Fatal("should not get here")
}

func TestSetAtomic(t *testing.T) {
	slice := make([]byte, 4)
	SetAtomic(slice, 31, true)
	if !Get(slice, 31) {
		t.Fatal("should return true")
	}
}
