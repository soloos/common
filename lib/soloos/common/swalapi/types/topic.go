package types

import (
	"bytes"
	sdbapitypes "soloos/common/sdbapi/types"
	"soloos/sdbone/offheap"
	"unsafe"
)

const (
	TopicStructSize  = unsafe.Sizeof(Topic{})
	TopicIDBytesNums = 64
)

type TopicID struct {
	Bytes    [TopicIDBytesNums]byte
	bytesLen int
}
type TopicUintptr uintptr

func StrToTopicID(topicIDStr string) TopicID {
	var ret TopicID
	ret.SetStr(topicIDStr)
	return ret
}

func (p TopicID) Str() string {
	return string(p.Bytes[:p.bytesLen])
}

func (p *TopicID) SetStr(topicIDStr string) {
	p.bytesLen = len(topicIDStr)
	copy(p.Bytes[:p.bytesLen], topicIDStr)
}

func (p *TopicID) SetBytes(in [TopicIDBytesNums]byte) {
	var index = bytes.IndexByte(in[:], 0)
	if index == -1 {
		p.bytesLen = TopicIDBytesNums
		copy(p.Bytes[:], in[:])
	} else {
		p.bytesLen = index
		copy(p.Bytes[:index], in[:index])
	}
}

func (u TopicUintptr) Ptr() *Topic { return (*Topic)(unsafe.Pointer(u)) }

type TopicMeta struct {
	TopicID         TopicID
	SWALMemberGroup SWALMemberGroup
}

type Topic struct {
	offheap.LKVTableObjectWithBytes64
	IsDBMetaDataInited sdbapitypes.MetaDataState
	Meta               TopicMeta
}

func (p *Topic) Reset() {
	p.IsDBMetaDataInited.Reset()
}
