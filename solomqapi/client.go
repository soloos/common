package solomqapi

import "soloos/common/solofsapitypes"

type SendLogSig = bool

type Client interface {
	SetSolofsClient(itSolofsClient interface{}) error
	PrepareNetBlockMetaData(uNetBlock solofsapitypes.NetBlockUintptr,
		uNetINode solofsapitypes.NetINodeUintptr, netblockIndex int32) error
	UploadMemBlockWithSolomq(uJob solofsapitypes.UploadMemBlockJobUintptr,
		uploadPeerIndex int) error
}
