package solofstypes

import (
	"soloos/common/solodbtypes"
	"soloos/common/soloosbase"
	"soloos/common/util"
	"unsafe"
)

const (
	UploadMemBlockJobStructSize = unsafe.Sizeof(UploadMemBlockJob{})
)

type UploadMemBlockJobID = int64

type UploadMemBlockJobUintptr uintptr

func (u UploadMemBlockJobUintptr) Ptr() *UploadMemBlockJob {
	return (*UploadMemBlockJob)(unsafe.Pointer(u))
}

type UploadMemBlockJob struct {
	MetaDataState solodbtypes.MetaDataState
	SyncDataSig   util.WaitGroup
	UNetINode     NetINodeUintptr
	UNetBlock     NetBlockUintptr
	UMemBlock     MemBlockUintptr
	MemBlockIndex MemBlockIndex
	// ID            UploadMemBlockJobID
	// NetINodeID    NetINodeID
	// NetBlockIndex NetBlockIndex
	// MemBlockIndex MemBlockIndex
	soloosbase.UploadBlockJob
}

func (p *UploadMemBlockJob) Reset() {
	p.MetaDataState.Store(solodbtypes.MetaDataStateUninited)
}
