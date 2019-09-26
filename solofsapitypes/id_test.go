package solofsapitypes

import (
	"soloos/common/soloosbase"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func BenchmarkEncodeBindIndex(b *testing.B) {
	var (
		u     uintptr = 0x12
		index int32   = 3
		id    soloosbase.PtrBindIndex
	)
	for n := 0; n < b.N; n++ {
		soloosbase.EncodePtrBindIndex(&id, u, index)
	}
}

func TestEncodeNetINodeBlockID(t *testing.T) {
	var netINodeBlockID NetINodeBlockID
	netINodeID := NetINodeID{1, 2, 3}
	blockIndex := int32(22)
	EncodeNetINodeBlockID(&netINodeBlockID, netINodeID, blockIndex)
	assert.Equal(t, uint8(1), netINodeBlockID[0])
	assert.Equal(t, uint8(2), netINodeBlockID[1])
	assert.Equal(t, uint8(3), netINodeBlockID[2])
	assert.Equal(t, uint8(22), netINodeBlockID[NetINodeIDSize])
}

func TestEncodePtrBindIndex(t *testing.T) {
	var (
		u     uintptr = 0x12
		index int32   = 3
		id    soloosbase.PtrBindIndex
	)
	soloosbase.EncodePtrBindIndex(&id, u, index)
	assert.Equal(t, uintptr(0x12), *((*uintptr)(unsafe.Pointer(&id))))
	assert.Equal(t, uint8(3), id[UintptrSize])
}

func BenchmarkEncodeNetINodeBlockID(b *testing.B) {
	var netINodeBlockID NetINodeBlockID
	netINodeID := NetINodeID{1, 2, 3}
	blockIndex := int32(22)
	for n := 0; n < b.N; n++ {
		EncodeNetINodeBlockID(&netINodeBlockID, netINodeID, blockIndex)
	}
}
