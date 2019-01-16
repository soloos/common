package offheap

import (
	"sync"
	"sync/atomic"
)

type SharedPointerBase struct {
	accessRWMutex sync.RWMutex
	Accessor      int32
}

func (p *SharedPointerBase) ReadAcquire() {
	atomic.AddInt32(&p.Accessor, 1)
	p.accessRWMutex.RLock()
}

func (p *SharedPointerBase) ReadRelease() {
	p.accessRWMutex.RUnlock()
	atomic.AddInt32(&p.Accessor, -1)
}

func (p *SharedPointerBase) WriteAcquire() {
	atomic.AddInt32(&p.Accessor, 1)
	p.accessRWMutex.Lock()
}

func (p *SharedPointerBase) WriteRelease() {
	p.accessRWMutex.Unlock()
	atomic.AddInt32(&p.Accessor, -1)
}

type SharedPointer struct {
	SharedPointerBase
	IsInited        bool
	IsShouldRelease bool
}

func (p *SharedPointer) SetReleasable() {
	p.IsShouldRelease = true
}

func (p *SharedPointer) Reset() {
	p.IsInited = false
	p.IsShouldRelease = false
}
