package offheap

import "sync/atomic"

func (p *OffheapDriver) AllocRawChunkPoolID() int32 {
	return atomic.AddInt32(&p.maxRawChunkPoolID, 1)
}

func (p *OffheapDriver) InitRawObjectPool(pool *RawObjectPool,
	structSize int, rawChunksLimit int32,
	prepareNewRawChunkFunc RawChunkPoolInvokePrepareNewRawChunk,
	releaseRawChunkFunc RawChunkPoolInvokeReleaseRawChunk) error {
	var (
		err error
	)

	err = pool.Init(p.AllocRawChunkPoolID(), structSize, rawChunksLimit, prepareNewRawChunkFunc, releaseRawChunkFunc)
	if err != nil {
		return err
	}

	p.SetRawChunkPool(&pool.memRawChunkPool)

	return nil
}

func (p *OffheapDriver) SetRawChunkPool(rawChunkPool *RawChunkPool) {
	p.rawChunkPools[rawChunkPool.ID] = rawChunkPool
}

func (p *OffheapDriver) GetRawChunkPool(poolid int32) *RawChunkPool {
	return p.rawChunkPools[poolid]
}
