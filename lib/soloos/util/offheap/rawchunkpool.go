package offheap

import (
	"math"
	"sync"
	"sync/atomic"
)

type RawChunkPoolInvokePrepareNewRawChunk func(uRawChunk uintptr)
type RawChunkPoolInvokeReleaseRawChunk func()

// RawChunkPool
// user -> AllocRawChunk -> mallocRawChunk -> user
// user -> AllocRawChunk -> RawChunkPoolAssistant.RawChunkPoolInvokeReleaseRawChunk -> ReleaseRawChunk -> user
type RawChunkPool struct {
	ID int32

	options RawChunkPoolOptions

	perRawChunkSize  uintptr
	perMmapBytesSize int
	currentMmapBytes *mmapbytes
	mmapBytesList    []*mmapbytes

	maxRawChunkID      int32
	chunksMutex        sync.Mutex
	activeRawChunksNum int32
	pool               sync.NoGCUintptrPool
	chunks             map[uintptr]uintptr
}

func (p *RawChunkPool) Init(id int32, options RawChunkPoolOptions) error {
	var (
		err error
	)

	p.ID = id
	p.options = options
	p.perRawChunkSize = uintptr(p.options.RawChunkSize)
	if p.options.RawChunksLimit == -1 {
		p.perMmapBytesSize = int(1024 * int(p.perRawChunkSize))
	} else {
		p.perMmapBytesSize = int(math.Ceil(float64(p.options.RawChunksLimit)/float64(16))) * int(p.perRawChunkSize)
	}

	err = p.grawMmapBytesList()
	if err != nil {
		return err
	}

	p.activeRawChunksNum = 0
	p.pool.New = p.mallocRawChunk

	return nil
}

func (p *RawChunkPool) mallocRawChunk() uintptr {
	var (
		uRawChunk        uintptr
		currentMmapBytes *mmapbytes
		end              uintptr
		err              error
	)

	// step1 graw mem if need
	if err == nil {
		currentMmapBytes = p.currentMmapBytes
		end = atomic.AddUintptr(&currentMmapBytes.readOff, p.perRawChunkSize)
		if end > currentMmapBytes.readBoundary {
			p.chunksMutex.Lock()
			currentMmapBytes = p.currentMmapBytes
			end = atomic.AddUintptr(&currentMmapBytes.readOff, p.perRawChunkSize)
			if end < currentMmapBytes.readBoundary {
				p.chunksMutex.Unlock()
				goto STEP1_DONE
			}

			err = p.grawMmapBytesList()
			if err != nil {
				p.chunksMutex.Unlock()
				goto STEP1_DONE
			}

			currentMmapBytes = p.currentMmapBytes
			end = currentMmapBytes.readOff + p.perRawChunkSize
			currentMmapBytes.readOff = end
			p.chunksMutex.Unlock()
		}
	}
STEP1_DONE:

	// step2 alloc mem for chunk
	if err == nil {
		// get chunk address
		uRawChunk = (uintptr)(end - p.perRawChunkSize)
	}

	if err != nil {
		panic("malloc chunk error")
	}

	p.options.RawChunkPoolInvokePrepareNewRawChunk(uRawChunk)
	return uintptr(uRawChunk)
}

func (p *RawChunkPool) AllocRawChunk() uintptr {
	if p.options.RawChunksLimit == -1 {
		return p.pool.Get()
	}

	// assert p.options.RawChunkReleaser != nil
	if atomic.AddInt32(&p.activeRawChunksNum, 1) < p.options.RawChunksLimit {
		return p.pool.Get()
	}

	for p.activeRawChunksNum > p.options.RawChunksLimit {
		p.options.RawChunkPoolInvokeReleaseRawChunk()
	}

	return p.pool.Get()
}

func (p *RawChunkPool) ReleaseRawChunk(chunk uintptr) {
	atomic.AddInt32(&p.activeRawChunksNum, -1)
	p.pool.Put(uintptr(chunk))
}
