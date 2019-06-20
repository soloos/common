package sdfsapitypes

import (
	"soloos/common/sdbapitypes"
	"soloos/common/soloosbase"
	"soloos/common/util"
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
	SyncDataSig   util.WaitGroup
	UNetINode     NetINodeUintptr
	UNetBlock     NetBlockUintptr
	UMemBlock     MemBlockUintptr
	MemBlockIndex int32
	soloosbase.UploadBlockJob
}

func (p *UploadMemBlockJob) Reset() {
	p.MetaDataState.Store(sdbapitypes.MetaDataStateUninited)
}
