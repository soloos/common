package offheap

import "sync/atomic"

func (p *OffheapDriver) AllocRawChunkPoolID() int32 {
	return atomic.AddInt32(&p.maxRawChunkPoolID, 1)
}

func (p *OffheapDriver) SetRawChunkPool(rawChunkPool *RawChunkPool) {
	p.rawChunkPools[rawChunkPool.ID] = rawChunkPool
}

func (p *OffheapDriver) GetRawChunkPool(poolid int32) *RawChunkPool {
	return p.rawChunkPools[poolid]
}
