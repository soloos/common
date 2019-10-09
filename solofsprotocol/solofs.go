package solofsprotocol

import (
	"encoding/gob"
	"soloos/common/snet"
	"soloos/common/solofstypes"
)

//go:generate msgp

func init() {
	gob.Register(SNetPeer{})
	gob.Register(NetINodeCommitSizeInDBReq{})
	gob.Register(NetINodeInfoReq{})
	gob.Register(NetINodeInfoResp{})
	gob.Register(NetINodeNetBlockInfoReq{})
	gob.Register(NetINodeNetBlockInfoResp{})
	gob.Register(NetINodePWriteReq{})
	gob.Register(NetINodeSyncReq{})
	gob.Register(NetINodePReadReq{})
	gob.Register(NetINodePReadResp{})
}

type SNetPeer struct {
	PeerID   string
	Address  string
	Protocol string
}

type NetINodeCommitSizeInDBReq struct {
	NetINodeID solofstypes.NetINodeID
	Size       uint64
}

type NetINodeInfoReq struct {
	NetINodeID  solofstypes.NetINodeID
	Size        uint64
	NetBlockCap int32
	MemBlockCap int32
}

type NetINodeInfoResp struct {
	snet.RespCommon
	Size        uint64
	NetBlockCap int32
	MemBlockCap int32
}

type NetINodeNetBlockInfoReq struct {
	NetINodeID    solofstypes.NetINodeID
	NetBlockIndex int32
	Cap           int32
}

type NetINodeNetBlockInfoResp struct {
	snet.RespCommon
	Len      int32
	Cap      int32
	Backends []string
}

type NetINodePWriteReq struct {
	NetINodeID       solofstypes.NetINodeID
	Offset           uint64
	Length           int32
	TransferBackends []string
}

type NetINodeSyncReq struct {
	NetINodeID solofstypes.NetINodeID
}

type NetINodePReadReq struct {
	NetINodeID solofstypes.NetINodeID
	Offset     uint64
	Length     int32
}

type NetINodePReadResp struct {
	snet.RespCommon
	Length int32
}
