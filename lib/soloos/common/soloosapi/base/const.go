package base

import "unsafe"

const (
	UUintptrSize = unsafe.Sizeof(uintptr(0))
	UintptrSize  = int(unsafe.Sizeof(uintptr(0)))
	Int32Size    = int(unsafe.Sizeof(int32(0)))
)
