package sdbapitypes

import (
	"sync"
	"sync/atomic"
)

const (
	MetaDataStateUninited int32 = iota
	MetaDataStateInited
)

type MetaDataState struct {
	ContextMutex sync.Mutex
	State        int32
}

func (p *MetaDataState) Load() int32 {
	return atomic.LoadInt32(&p.State)
}

func (p *MetaDataState) Store(v int32) {
	atomic.StoreInt32(&p.State, v)
}

func (p *MetaDataState) Reset() {
	atomic.StoreInt32(&p.State, MetaDataStateUninited)
}

func (p *MetaDataState) LockContext() {
	p.ContextMutex.Lock()
}

func (p *MetaDataState) UnlockContext() {
	p.ContextMutex.Unlock()
}
