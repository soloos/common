package base

import (
	"soloos/sdbone/offheap"
	"sync"
	"sync/atomic"
)

type UploadBlockJob struct {
	uploadMaskMutex           sync.Mutex
	uploadMasks               [2]offheap.ChunkMask
	uploadMaskProcessingIndex int32
}

func (p *UploadBlockJob) swapUploadMask() {
	if atomic.CompareAndSwapInt32(&p.uploadMaskProcessingIndex, 0, 1) == false {
		atomic.CompareAndSwapInt32(&p.uploadMaskProcessingIndex, 1, 0)
	}
}

func (p *UploadBlockJob) getUploadMaskWaitingIndex() int32 {
	if atomic.LoadInt32(&p.uploadMaskProcessingIndex) == 0 {
		return 1
	}
	return 0
}

func (p *UploadBlockJob) Reset() {
	p.uploadMasks[0].Reset()
	p.uploadMasks[1].Reset()
}

// PrepareUpload return is job done
func (p *UploadBlockJob) PrepareUploadMask() bool {
	// prepare upload job
	p.uploadMaskMutex.Lock()
	if p.uploadMasks[p.getUploadMaskWaitingIndex()].MaskArrayLen == 0 {
		// upload done and continue
		p.uploadMaskMutex.Unlock()
		return true
	}

	// upload job exists
	p.swapUploadMask()
	p.uploadMaskMutex.Unlock()
	return false
}

func (p *UploadBlockJob) WaitingQueueMergeIncludeNeighbour(offset, end int) (isMergeEventHappened, isMergeWriteMaskSuccess bool) {
	p.uploadMaskMutex.Lock()
	isMergeEventHappened, isMergeWriteMaskSuccess =
		p.uploadMasks[p.getUploadMaskWaitingIndex()].MergeIncludeNeighbour(offset, end)
	p.uploadMaskMutex.Unlock()
	return
}

func (p *UploadBlockJob) GetProcessingChunkMask() offheap.ChunkMask {
	return p.uploadMasks[p.uploadMaskProcessingIndex]
}

func (p *UploadBlockJob) ResetProcessingChunkMask() {
	p.uploadMasks[p.uploadMaskProcessingIndex].Reset()
}
