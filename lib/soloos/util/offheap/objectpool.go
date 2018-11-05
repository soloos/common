package offheap

import (
	"sync"
)

type ObjectPool struct {
	memChunkPool ChunkPool
	Objects      sync.Map
}

func (p *ObjectPool) Init(id int32, structSize int, chunksLimit int32,
	prepareNewChunkFunc ChunkPoolInvokePrepareNewChunk,
	releaseChunkFunc ChunkPoolInvokeReleaseChunk) error {
	var (
		err     error
		options ChunkPoolOptions
	)

	options.ChunkSize = structSize
	options.ChunksLimit = chunksLimit
	options.SetChunkPoolAssistant(prepareNewChunkFunc, releaseChunkFunc)
	err = p.memChunkPool.Init(id, options)
	if err != nil {
		return err
	}

	return nil
}

func (p *ObjectPool) AllocObject() uintptr {
	return p.memChunkPool.AllocChunk().Ptr().Data
}

func (p *ObjectPool) ReleaseObjectByID(id interface{}) {
	retI, exists := p.Objects.Load(id)
	if exists {
		p.ReleaseObject(retI.(uintptr))
		p.Objects.Delete(id)
	}
}
func (p *ObjectPool) ReleaseObject(uObject uintptr) {
	p.memChunkPool.ReleaseChunk(ChunkUintptr(uObject - ChunkStructSize))
}

// MustGetChunk get or init a filechunk
// The last result is true if the value was loaded, false if alloc.
func (p *ObjectPool) MustGetObject(id interface{}) (uintptr, bool) {
	var (
		retI    interface{}
		uObject uintptr
		exists  bool
	)
	retI, exists = p.Objects.Load(id)
	if exists {
		return retI.(uintptr), true
	}

	uObject = p.AllocObject()
	retI, exists = p.Objects.LoadOrStore(id, uObject)
	if exists {
		p.ReleaseObject(uObject)
		return retI.(uintptr), true
	}

	return uObject, false
}
