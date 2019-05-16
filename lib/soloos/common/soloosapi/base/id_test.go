package base

import "testing"

func BenchmarkEncodeBindIndex(b *testing.B) {
	var (
		u     uintptr = 0x12
		index int32   = 3
		id    PtrBindIndex
	)
	for n := 0; n < b.N; n++ {
		soloosbase.EncodePtrBindIndex(&id, u, index)
	}
}
