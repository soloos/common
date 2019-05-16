package types

import (
	sdbapitypes "soloos/common/sdbapi/types"
	soloosbase "soloos/common/soloosapi/base"
	"sync"
	"unsafe"
)

const (
	UploadMemBlockJobStructSize = unsafe.Sizeof(UploadMemBlockJob{})
)

type UploadMemBlockJobUintptr uintptr

func (u UploadMemBlockJobUintptr) Ptr() *UploadMemBlockJob {
	return (*UploadMemBlockJob)(unsafe.Pointer(u))
}

type UploadMemBlockJob struct {
	MetaDataState sdbapitypes.MetaDataState
	SyncDataSig   sync.WaitGroup
	UNetINode     NetINodeUintptr
	UNetBlock     NetBlockUintptr
	UMemBlock     MemBlockUintptr
	MemBlockIndex int32
	soloosbase.UploadBlockJob
}

func (p *UploadMemBlockJob) Reset() {
	p.MetaDataState.Store(sdbapitypes.MetaDataStateUninited)
}
