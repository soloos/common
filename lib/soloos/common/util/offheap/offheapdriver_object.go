package offheap

func (p *OffheapDriver) InitObjectPool(pool *ObjectPool,
	structSize int, chunksLimit int32,
	prepareNewChunkFunc ChunkPoolInvokePrepareNewChunk,
	releaseChunkFunc ChunkPoolInvokeReleaseChunk) error {
	var (
		err error
	)

	err = pool.Init(p.AllocChunkPoolID(), structSize, chunksLimit, prepareNewChunkFunc, releaseChunkFunc)
	if err != nil {
		return err
	}

	p.SetChunkPool(&pool.chunkPool)

	return nil
}
