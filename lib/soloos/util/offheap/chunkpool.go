package offheap

import (
	"math"
	"sync"
	"sync/atomic"
)

type ChunkPoolInvokePrepareNewChunk func(uChunk ChunkUintptr)
type ChunkPoolInvokeReleaseChunk func()

// ChunkPool
// user -> AllocChunk -> mallocChunk -> user
// user -> AllocChunk -> ChunkPoolAssistant.ChunkPoolInvokeReleaseChunk -> ReleaseChunk -> user
type ChunkPool struct {
	ID int32

	options ChunkPoolOptions

	perChunkSize           uintptr
	perChunkWithStructSize uintptr
	perMmapBytesSize       int
	currentMmapBytes       *mmapbytes
	mmapBytesList          []*mmapbytes

	maxChunkID      int32
	chunksMutex     sync.Mutex
	activeChunksNum int32
	pool            sync.NoGCUintptrPool
	chunks          map[uintptr]uintptr
}

func (p *ChunkPool) Init(id int32, options ChunkPoolOptions) error {
	var (
		err error
	)

	p.ID = id
	p.options = options
	p.perChunkSize = uintptr(p.options.ChunkSize)
	p.perChunkWithStructSize = ChunkStructSize + p.perChunkSize
	if p.options.ChunksLimit == -1 {
		p.perMmapBytesSize = int(1024 * int(p.perChunkWithStructSize))
	} else {
		p.perMmapBytesSize = int(math.Ceil(float64(p.options.ChunksLimit)/float64(16))) * int(p.perChunkWithStructSize)
	}

	err = p.growMmapBytesList()
	if err != nil {
		return err
	}

	p.activeChunksNum = 0
	p.pool.New = p.mallocChunk

	return nil
}

func (p *ChunkPool) mallocChunk() uintptr {
	var (
		uChunk           ChunkUintptr
		currentMmapBytes *mmapbytes
		end              uintptr
		err              error
	)

	// step1 grow mem if need
	if err == nil {
		currentMmapBytes = p.currentMmapBytes
		end = atomic.AddUintptr(&currentMmapBytes.readOff, p.perChunkWithStructSize)
		if end > currentMmapBytes.readBoundary {
			p.chunksMutex.Lock()
			currentMmapBytes = p.currentMmapBytes
			end = atomic.AddUintptr(&currentMmapBytes.readOff, p.perChunkWithStructSize)
			if end < currentMmapBytes.readBoundary {
				p.chunksMutex.Unlock()
				goto STEP1_DONE
			}

			err = p.growMmapBytesList()
			if err != nil {
				p.chunksMutex.Unlock()
				goto STEP1_DONE
			}

			currentMmapBytes = p.currentMmapBytes
			end = currentMmapBytes.readOff + p.perChunkWithStructSize
			currentMmapBytes.readOff = end
			p.chunksMutex.Unlock()
		}
	}
STEP1_DONE:

	// step2 alloc mem for chunk
	if err == nil {
		// get chunk address
		uChunk = (ChunkUintptr)(end - p.perChunkWithStructSize)

		// save chunk in offheap.data
		uChunk.Ptr().ID = atomic.AddInt32(&p.maxChunkID, 1)
		uChunk.Ptr().Data = uintptr(uChunk) + ChunkStructSize
	}

	if err != nil {
		panic("malloc chunk error")
	}

	p.options.ChunkPoolInvokePrepareNewChunk(uChunk)
	return uintptr(uChunk)
}

func (p *ChunkPool) AllocChunk() ChunkUintptr {
	if p.options.ChunksLimit == -1 {
		return ChunkUintptr(p.pool.Get())
	}

	// assert p.options.ChunkReleaser != nil
	if atomic.AddInt32(&p.activeChunksNum, 1) < p.options.ChunksLimit {
		return ChunkUintptr(p.pool.Get())
	}

	for p.activeChunksNum > p.options.ChunksLimit {
		p.options.ChunkPoolInvokeReleaseChunk()
	}

	return ChunkUintptr(p.pool.Get())
}

func (p *ChunkPool) ReleaseChunk(chunk ChunkUintptr) {
	atomic.AddInt32(&p.activeChunksNum, -1)
	p.pool.Put(uintptr(chunk))
}
