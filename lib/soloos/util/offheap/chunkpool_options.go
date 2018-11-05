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
