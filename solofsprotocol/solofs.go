package solofsprotocol

import (
	"encoding/gob"
	"soloos/common/snet"
	"soloos/common/solofsapitypes"
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
	NetINodeID solofsapitypes.NetINodeID
	Size       uint64
}

type NetINodeInfoReq struct {
	NetINodeID  solofsapitypes.NetINodeID
	Size        uint64
	NetBlockCap int32
	MemBlockCap int32
}

type NetINodeInfoResp struct {
	snet.RespDataCommon
	Size        uint64
	NetBlockCap int32
	MemBlockCap int32
}

var _ = snet.IRespData(&NetINodeInfoResp{})

type NetINodeNetBlockInfoReq struct {
	NetINodeID    solofsapitypes.NetINodeID
	NetBlockIndex int32
	Cap           int32
}

type NetINodeNetBlockInfoResp struct {
	snet.RespDataCommon
	Len      int32
	Cap      int32
	Backends []string
}

var _ = snet.IRespData(&NetINodeNetBlockInfoResp{})

type NetINodePWriteReq struct {
	NetINodeID       solofsapitypes.NetINodeID
	Offset           uint64
	Length           int32
	TransferBackends []string
}

type NetINodeSyncReq struct {
	NetINodeID solofsapitypes.NetINodeID
}

type NetINodePReadReq struct {
	NetINodeID solofsapitypes.NetINodeID
	Offset     uint64
	Length     int32
}

type NetINodePReadResp struct {
	snet.RespDataCommon
	Length int32
}

var _ = snet.IRespData(&NetINodePReadResp{})
