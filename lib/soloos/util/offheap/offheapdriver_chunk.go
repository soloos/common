package offheap

import "sync/atomic"

func (p *OffheapDriver) AllocChunkPoolID() int32 {
	return atomic.AddInt32(&p.maxChunkPoolID, 1)
}

func (p *OffheapDriver) InitChunkPool(options ChunkPoolOptions, chunkPool *ChunkPool) error {
	err := chunkPool.Init(p.AllocChunkPoolID(), options)
	if err != nil {
		return err
	}

	p.chunkPools[chunkPool.ID] = chunkPool
	return nil
}

func (p *OffheapDriver) InitObjectPool(pool *ObjectPool,
	structSize int, rawChunksLimit int32,
	prepareNewChunkFunc ChunkPoolInvokePrepareNewChunk,
	releaseChunkFunc ChunkPoolInvokeReleaseChunk) error {
	var (
		err error
	)

	err = pool.Init(p.AllocChunkPoolID(), structSize, rawChunksLimit, prepareNewChunkFunc, releaseChunkFunc)
	if err != nil {
		return err
	}

	p.SetChunkPool(&pool.memChunkPool)

	return nil
}

func (p *OffheapDriver) SetChunkPool(chunkPool *ChunkPool) {
	p.chunkPools[chunkPool.ID] = chunkPool
}

func (p *OffheapDriver) GetChunkPool(poolid int32) *ChunkPool {
	return p.chunkPools[poolid]
}
