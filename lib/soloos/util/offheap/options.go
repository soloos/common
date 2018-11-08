package offheap

type ChunkPoolOptions struct {
	ChunkSize   int
	ChunksLimit int32

	ChunkPoolInvokePrepareNewChunk
	ChunkPoolInvokeReleaseChunk
}

func MakeDefaultTestChunkPoolOptions(chunkSize int) ChunkPoolOptions {
	return ChunkPoolOptions{
		ChunkSize:   chunkSize,
		ChunksLimit: 1024,
	}
}

func (p *ChunkPoolOptions) SetChunkPoolAssistant(prepareNewChunkFunc ChunkPoolInvokePrepareNewChunk,
	releaseChunkFunc ChunkPoolInvokeReleaseChunk) {
	p.ChunkPoolInvokePrepareNewChunk = prepareNewChunkFunc
	p.ChunkPoolInvokeReleaseChunk = releaseChunkFunc
}

type RawChunkPoolOptions struct {
	RawChunkSize   int
	RawChunksLimit int32

	RawChunkPoolInvokePrepareNewRawChunk
	RawChunkPoolInvokeReleaseRawChunk
}

func MakeDefaultTestRawChunkPoolOptions(chunkSize int) RawChunkPoolOptions {
	return RawChunkPoolOptions{
		RawChunkSize:   chunkSize,
		RawChunksLimit: 1024,
	}
}

func (p *RawChunkPoolOptions) SetRawChunkPoolAssistant(prepareNewRawChunkFunc RawChunkPoolInvokePrepareNewRawChunk,
	releaseRawChunkFunc RawChunkPoolInvokeReleaseRawChunk) {
	p.RawChunkPoolInvokePrepareNewRawChunk = prepareNewRawChunkFunc
	p.RawChunkPoolInvokeReleaseRawChunk = releaseRawChunkFunc
}
