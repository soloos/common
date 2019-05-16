package types

import (
	soloosbase "soloos/common/soloosapi/base"
	"unsafe"
)

const (
	UploadTopicMsgJobStructSize = unsafe.Sizeof(UploadTopicMsgJob{})
)

type UploadTopicMsgJobUintptr uintptr

func (u UploadTopicMsgJobUintptr) Ptr() *UploadTopicMsgJob {
	return (*UploadTopicMsgJob)(unsafe.Pointer(u))
}

type UploadTopicMsgJob struct {
	TopicID       TopicID
	NetBlockID    [64]byte
	MemBlockIndex int
	soloosbase.UploadBlockJob
}
