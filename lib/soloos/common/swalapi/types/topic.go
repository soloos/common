package types

import (
	"soloos/sdbone/offheap"
	"unsafe"
)

const (
	TopicStructSize  = unsafe.Sizeof(Topic{})
	TopicIDBytesNums = 64
)

type TopicID struct {
	Data    [TopicIDBytesNums]byte
	dataLen int
}
type TopicUintptr uintptr

func StrToTopicID(topicIDStr string) TopicID {
	var ret TopicID
	ret.SetStr(topicIDStr)
	return ret
}

func (p TopicID) Str() string {
	return string(p.Data[:p.dataLen])
}

func (p *TopicID) SetStr(topicIDStr string) {
	p.dataLen = len(topicIDStr)
	copy(p.Data[:p.dataLen], topicIDStr)
}

type Topic struct {
	offheap.LKVTableObjectWithBytes64
	TopicID         TopicID
	SWALMemberGroup SWALMemberGroup
}
