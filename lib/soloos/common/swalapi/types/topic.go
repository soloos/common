package types

import "unsafe"

const (
	TopicStructSize = unsafe.Sizeof(Topic{})
	TopicIDBytesCap = 64
)

type TopicID struct {
	dataLen int
	Data    [TopicIDBytesCap]byte
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
	TopicID         TopicID
	SWALMemberGroup SWALMemberGroup
}
