package util

import (
	"sync/atomic"

	deadlock "github.com/sasha-s/go-deadlock"
)

// type RawWaitGroup = sync.WaitGroup
type RawWaitGroup = WaitGroup

type WaitGroup struct {
	waitProtect deadlock.Mutex
	deadlock.WaitGroup

	AccessorNum int32
}

func (p *WaitGroup) Add(delta int) {
	atomic.AddInt32(&p.AccessorNum, int32(delta))

	p.WaitGroup.Add(delta)
}

func (p *WaitGroup) Done() {
	atomic.AddInt32(&p.AccessorNum, -1)
	p.WaitGroup.Done()
}

func (p *WaitGroup) Wait() {
	p.waitProtect.Lock()
	p.WaitGroup.Wait()
	p.waitProtect.Unlock()
}

type Mutex = deadlock.Mutex
type RWMutex = deadlock.RWMutex

func init() {
}
