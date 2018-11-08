package offheap

type OffheapDriver struct {
	chunkPools     map[int32]*ChunkPool
	maxChunkPoolID int32

	rawChunkPools     map[int32]*RawChunkPool
	maxRawChunkPoolID int32
}

func (p *OffheapDriver) Init() error {
	p.chunkPools = make(map[int32]*ChunkPool)
	p.rawChunkPools = make(map[int32]*RawChunkPool)
	return nil
}
