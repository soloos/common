package offheap

import "sync/atomic"

func (p *OffheapDriver) AllocChunkPoolID() int32 {
	return atomic.AddInt32(&p.maxChunkPoolID, 1)
}

func (p *OffheapDriver) InitChunkPool(pool *ChunkPool,
	chunkSize int, chunksLimit int32,
	prepareNewChunkFunc ChunkPoolInvokePrepareNewChunk,
	releaseChunkFunc ChunkPoolInvokeReleaseChunk) error {
	err := pool.Init(p.AllocChunkPoolID(), chunkSize, chunksLimit, prepareNewChunkFunc, releaseChunkFunc)
	if err != nil {
		return err
	}

	p.SetChunkPool(pool)
	return nil
}

func (p *OffheapDriver) SetChunkPool(chunkPool *ChunkPool) {
	p.chunkPools[chunkPool.ID] = chunkPool
}

func (p *OffheapDriver) GetChunkPool(poolid int32) *ChunkPool {
	return p.chunkPools[poolid]
}
