package offheap

var (
	DefaultOffheapDriver OffheapDriver
)

func init() {
	var err error
	err = DefaultOffheapDriver.Init()
	if err != nil {
		panic(err)
	}
}

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

func InitRawObjectPool(pool *RawObjectPool,
	structSize int, rawChunksLimit int32,
	prepareNewRawChunkFunc RawChunkPoolInvokePrepareNewRawChunk,
	releaseRawChunkFunc RawChunkPoolInvokeReleaseRawChunk) error {

	return DefaultOffheapDriver.InitRawObjectPool(pool,
		structSize, rawChunksLimit,
		prepareNewRawChunkFunc,
		releaseRawChunkFunc)
}

func InitObjectPool(pool *ObjectPool,
	structSize int, chunksLimit int32,
	prepareNewChunkFunc ChunkPoolInvokePrepareNewChunk,
	releaseChunkFunc ChunkPoolInvokeReleaseChunk) error {

	return DefaultOffheapDriver.InitObjectPool(pool,
		structSize, chunksLimit,
		prepareNewChunkFunc,
		releaseChunkFunc)
}

func InitChunkPool(pool *ChunkPool, chunkSize int, chunksLimit int32,
	prepareNewChunkFunc ChunkPoolInvokePrepareNewChunk,
	releaseChunkFunc ChunkPoolInvokeReleaseChunk) error {

	return DefaultOffheapDriver.InitChunkPool(pool,
		chunkSize, chunksLimit,
		prepareNewChunkFunc,
		releaseChunkFunc)
}
