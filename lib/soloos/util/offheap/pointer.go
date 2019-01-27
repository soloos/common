package offheap

import (
	"sync"
	"sync/atomic"
)

const (
	SharedPointerUninited   = int32(0)
	SharedPointerIniteded   = int32(1)
	SharedPointerReleasable = int32(2)
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
	Status int32
}

func (p *SharedPointer) SetReleasable() {
	atomic.StoreInt32(&p.Status, SharedPointerReleasable)
}

func (p *SharedPointer) Reset() {
	atomic.StoreInt32(&p.Status, SharedPointerUninited)
}

func (p *SharedPointer) CompleteInit() {
	atomic.StoreInt32(&p.Status, SharedPointerIniteded)
}

func (p *SharedPointer) IsInited() bool {
	return atomic.LoadInt32(&p.Status) > SharedPointerUninited
}

func (p *SharedPointer) IsShouldRelease() bool {
	return atomic.LoadInt32(&p.Status) == SharedPointerReleasable
}
