package offheap

import (
	"unsafe"
)

const (
	ChunkStructSize = unsafe.Sizeof(Chunk{})
)

type ChunkUintptr uintptr

func (p ChunkUintptr) Ptr() *Chunk { return (*Chunk)(unsafe.Pointer(p)) }

type Chunk struct {
	SharedPointerBase
	ID   int32
	Data uintptr
}
