package base

import "unsafe"

const (
	PtrBindIndexSize int = UintptrSize + Int32Size
)

type PtrBindIndex = [PtrBindIndexSize]byte

func EncodePtrBindIndex(id *PtrBindIndex, u uintptr, index int32) {
	*((*uintptr)(unsafe.Pointer(id))) = u
	*((*int32)(unsafe.Pointer(uintptr(unsafe.Pointer(id)) + UUintptrSize))) = index
}
