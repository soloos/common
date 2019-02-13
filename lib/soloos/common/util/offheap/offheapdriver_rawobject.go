package offheap

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

	p.SetRawChunkPool(&pool.rawChunkPool)

	return nil
}
