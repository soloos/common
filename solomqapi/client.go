package solomqapi

import "soloos/common/solofsapitypes"

type SendLogSig = bool

type Client interface {
	SetSOLOFSClient(itSOLOFSClient interface{}) error
	PrepareNetBlockMetaData(uNetBlock solofsapitypes.NetBlockUintptr,
		uNetINode solofsapitypes.NetINodeUintptr, netblockIndex int32) error
	UploadMemBlockWithSOLOMQ(uJob solofsapitypes.UploadMemBlockJobUintptr,
		uploadPeerIndex int) error
}
