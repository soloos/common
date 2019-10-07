package solomqprotocol

import (
	"encoding/gob"
	"soloos/common/solofstypes"
)

//go:generate msgp
func init() {
	gob.Register(TopicPrepareReq{})
	gob.Register(TopicPrepareNetBlockReq{})
	gob.Register(TopicPWriteReq{})
}

type TopicPrepareReq struct {
	TopicID   int64
	FsINodeID uint64
}

type TopicPrepareNetBlockReq struct {
	TopicID         int64
	FsINodeID       uint64
	IndexInNetINode int
}

type TopicPWriteReq struct {
	TopicID          int64
	NetINodeID       solofstypes.NetINodeID
	Offset           uint64
	Length           int
	TransferBackends []string
}
