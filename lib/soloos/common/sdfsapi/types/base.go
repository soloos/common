package types

import "sync/atomic"

type MetaDataState int32

func (p *MetaDataState) Load() int32 {
	return atomic.LoadInt32((*int32)(p))
}

func (p *MetaDataState) Store(v int32) {
	atomic.StoreInt32((*int32)(p), v)
}

func (p *MetaDataState) Reset() {
	atomic.StoreInt32((*int32)(p), MetaDataStateUninited)
}

const (
	MetaDataStateUninited int32 = iota
	MetaDataStateInited
)
