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
	rawObjectPool RawObjectPool
}

func (p *TPool) Init(id int32, structSize int, chunksLimit int32) {
	p.rawObjectPool.Init(id, structSize, chunksLimit,
		p.RawChunkPoolInvokePrepareNewRawChunk,
		p.RawChunkPoolInvokeReleaseRawChunk)
}

func (p *TPool) RawChunkPoolInvokePrepareNewRawChunk(uRawChunk uintptr) {
}

func (p *TPool) RawChunkPoolInvokeReleaseRawChunk() {
	p.rawObjectPool.RawObjects.Range(func(k, v interface{}) bool {
		p.rawObjectPool.ReleaseRawObject(p.rawObjectPool.ReleaseRawObjectByID(k))
		return true
	})
}

func BenchmarkRawObjectPool(b *testing.B) {
	runtime.GC()
	var tPool TPool
	tPool.Init(1, TStructSize, int32(b.N))

	for run := 0; run < 2; run++ {
		for n := 0; n < b.N; n++ {
			if n%10000 == 0 {
				runtime.GC()
			}
			tPool.rawObjectPool.ReleaseRawObject(tPool.rawObjectPool.AllocRawObject())
		}
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if n%10000 == 0 {
			runtime.GC()
		}
		tPool.rawObjectPool.AllocRawObject()
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

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if n%10000 == 0 {
			runtime.GC()
		}
		pool.Get()
	}
}

func BenchmarkRawObjectPoolMustGet(b *testing.B) {
	runtime.GC()
	debug.SetGCPercent(99)
	var tPool TPool
	tPool.Init(1, TStructSize, 102400)

	for n := 0; n < b.N; n++ {
		if n%100000 == 0 {
			runtime.GC()
		}
		tPool.rawObjectPool.MustGetRawObject(n)
	}
}

func TestRawObjectPool(t *testing.T) {
	var tPool TPool
	tPool.Init(1, TStructSize, 6)

	for n := 0; n < 10; n++ {
		tPool.rawObjectPool.MustGetRawObject(n)
	}
}
