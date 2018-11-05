package offheap

import (
	"runtime"
	"runtime/debug"
	"sync"
	"testing"
	"unsafe"
)

const TStructSize = int(unsafe.Sizeof(T{}))

type T struct {
	Data [1024]byte
}

type TPool struct {
	objectPool ObjectPool
}

func (p *TPool) Init(id int32, structSize int, chunksLimit int32) {
	p.objectPool.Init(id, structSize, chunksLimit, p.ChunkPoolInvokePrepareNewChunk, p.ChunkPoolInvokeReleaseChunk)
}

func (p *TPool) ChunkPoolInvokePrepareNewChunk(uChunk ChunkUintptr) {
}

func (p *TPool) ChunkPoolInvokeReleaseChunk() {
	p.objectPool.Objects.Range(func(k, v interface{}) bool {
		p.objectPool.ReleaseObjectByID(k)
		return true
	})
}

func TestObjectPool(t *testing.T) {
	var tPool TPool
	tPool.Init(1, TStructSize, 6)

	for n := 0; n < 10; n++ {
		tPool.objectPool.MustGetObject(n)
	}
}

func BenchmarkObjectPool(b *testing.B) {
	runtime.GC()
	debug.SetGCPercent(99)
	var tPool TPool
	tPool.Init(1, TStructSize, 102400)

	for n := 0; n < b.N; n++ {
		if n%100000 == 0 {
			runtime.GC()
		}
		tPool.objectPool.MustGetObject(n)
	}
}

func BenchmarkObjectPoolAlloc(b *testing.B) {
	runtime.GC()
	var tPool TPool
	tPool.Init(1, TStructSize, int32(b.N))

	for run := 0; run < 2; run++ {
		for n := 0; n < b.N; n++ {
			if n%10000 == 0 {
				runtime.GC()
			}
			tPool.objectPool.ReleaseObject(tPool.objectPool.AllocObject())
		}
	}

	for n := 0; n < b.N; n++ {
		if n%10000 == 0 {
			runtime.GC()
		}
		tPool.objectPool.AllocObject()
	}
}

func BenchmarkSyncPool(b *testing.B) {
	runtime.GC()
	var pool sync.Pool
	pool.New = func() interface{} {
		return new(T)
	}

	for run := 0; run < 2; run++ {
		for n := 0; n < b.N; n++ {
			if n%10000 == 0 {
				runtime.GC()
			}
			pool.Put(pool.Get())
		}
	}

	for n := 0; n < b.N; n++ {
		if n%10000 == 0 {
			runtime.GC()
		}
		pool.Get()
	}
}
