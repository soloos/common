package offheap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockChunkPool struct {
	driver        *MockOffheapDriver
	chunks        map[int32]ChunkUintptr
	offheapDriver *OffheapDriver
	chunkPool     ChunkPool
}

func (p *MockChunkPool) Init(chunks map[int32]ChunkUintptr, options ChunkPoolOptions, offheapDriver *OffheapDriver) error {
	var err error
	options.SetChunkPoolAssistant(p.ChunkPoolInvokePrepareNewChunk, p.ChunkPoolInvokeReleaseChunk)
	p.chunks = chunks
	p.offheapDriver = offheapDriver
	err = p.offheapDriver.InitChunkPool(options, &p.chunkPool)
	if err != nil {
		return err
	}

	return nil
}

func (p *MockChunkPool) takeChunkForRelease() ChunkUintptr {
	for _, uChunk := range p.chunks {
		return uChunk
	}
	return 0x00
}

func (p *MockChunkPool) ChunkPoolInvokePrepareNewChunk(uChunk ChunkUintptr) {
}

func (p *MockChunkPool) ChunkPoolInvokeReleaseChunk() {
	uChunk := p.takeChunkForRelease()
	pChunk := uChunk.Ptr()
	delete(p.chunks, pChunk.ID)
	p.chunkPool.ReleaseChunk(uChunk)
	return
}

func (p *MockChunkPool) AllocChunk() ChunkUintptr {
	uChunk := p.chunkPool.AllocChunk()
	p.chunks[uChunk.Ptr().ID] = uChunk
	return uChunk
}

func BenchmarkChunkPool(b *testing.B) {
	var (
		offheapDriver    OffheapDriver
		mockChunkPool    MockChunkPool
		chunkPoolOptions = MakeDefaultTestChunkPoolOptions(10)
	)

	offheapDriver.Init()
	mockChunkPool.Init(make(map[int32]ChunkUintptr), chunkPoolOptions, &offheapDriver)
	for n := 0; n < b.N; n++ {
		mockChunkPool.AllocChunk()
	}
}

func TestChunkPool(t *testing.T) {
	var (
		offheapDriver    OffheapDriver
		mockChunkPool    MockChunkPool
		chunkPoolOptions = MakeDefaultTestChunkPoolOptions(1024)
		uChunk           ChunkUintptr
	)

	assert.NoError(t, offheapDriver.Init())
	assert.NoError(t, mockChunkPool.Init(make(map[int32]ChunkUintptr), chunkPoolOptions, &offheapDriver))
	uChunk = mockChunkPool.chunkPool.AllocChunk()
	assert.NotNil(t, uChunk)

	mockChunkPool.chunkPool.ReleaseChunk(uChunk)
}
