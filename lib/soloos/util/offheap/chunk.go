package offheap

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

const (
	ChunkStructSize = unsafe.Sizeof(Chunk{})
)

type ChunkUintptr uintptr

func (p ChunkUintptr) Ptr() *Chunk { return (*Chunk)(unsafe.Pointer(p)) }

type Chunk struct {
	accessRWMutex sync.RWMutex
	Accessor      int32
	ID            int32
	Data          uintptr
}

func (p *Chunk) ReadAcquire() {
	atomic.AddInt32(&p.Accessor, 1)
	p.accessRWMutex.RLock()
}

func (p *Chunk) ReadRelease() {
	p.accessRWMutex.RUnlock()
	atomic.AddInt32(&p.Accessor, -1)
}

func (p *Chunk) WriteAcquire() {
	atomic.AddInt32(&p.Accessor, 1)
	p.accessRWMutex.Lock()
}

func (p *Chunk) WriteRelease() {
	p.accessRWMutex.Unlock()
	atomic.AddInt32(&p.Accessor, -1)
}
