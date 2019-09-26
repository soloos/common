package solomqapitypes

import (
	"bytes"
	"soloos/common/solodbapitypes"
	"soloos/solodb/offheap"
	"unsafe"
)

const (
	TopicStructSize    = unsafe.Sizeof(Topic{})
	TopicNameBytesNums = 64
)

type TopicName struct {
	Bytes    [TopicNameBytesNums]byte
	bytesLen int
}
type TopicID = int64
type TopicUintptr uintptr

func StrToTopicName(topicIDStr string) TopicName {
	var ret TopicName
	ret.SetStr(topicIDStr)
	return ret
}

func (p TopicName) Str() string {
	return string(p.Bytes[:p.bytesLen])
}

func (p *TopicName) SetStr(topicIDStr string) {
	p.bytesLen = len(topicIDStr)
	copy(p.Bytes[:p.bytesLen], topicIDStr)
}

func (p *TopicName) SetBytes(in [TopicNameBytesNums]byte) {
	var index = bytes.IndexByte(in[:], 0)
	if index == -1 {
		p.bytesLen = TopicNameBytesNums
		copy(p.Bytes[:], in[:])
	} else {
		p.bytesLen = index
		copy(p.Bytes[:index], in[:index])
	}
}

func (u TopicUintptr) Ptr() *Topic { return (*Topic)(unsafe.Pointer(u)) }

type TopicMeta struct {
	TopicID         TopicID
	TopicName       TopicName
	SOLOMQMemberGroup SOLOMQMemberGroup
}

type Topic struct {
	offheap.LKVTableObjectWithInt64
	Meta               TopicMeta
	IsDBMetaDataInited solodbapitypes.MetaDataState
}

func (p *Topic) Reset() {
	p.IsDBMetaDataInited.Reset()
}
