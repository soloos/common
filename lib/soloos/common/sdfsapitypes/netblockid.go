package sdfsapitypes

import (
	"reflect"
	"unsafe"
)

const (
	NetINodeBlockIDSize int = NetINodeIDSize + Int32Size
)

type NetINodeBlockID = [NetINodeBlockIDSize]byte

func EncodeNetINodeBlockID(netINodeBlockID *NetINodeBlockID, netINodeID NetINodeID, blockIndex int32) {
	bytes := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		uintptr(unsafe.Pointer(netINodeBlockID)),
		NetINodeBlockIDSize,
		NetINodeBlockIDSize,
	}))
	copy(bytes[:NetINodeIDSize], (*(*[NetINodeIDSize]byte)((unsafe.Pointer)(&netINodeID)))[:NetINodeIDSize])
	copy(bytes[NetINodeIDSize:], (*(*[Int32Size]byte)((unsafe.Pointer)(&blockIndex)))[:Int32Size])
}
