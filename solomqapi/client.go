package solomqapi

import "soloos/common/solofstypes"

type SendLogSig = bool

type Client interface {
	SetSolofsClient(itSolofsClient interface{}) error
	PrepareNetBlockMetaData(uNetBlock solofstypes.NetBlockUintptr,
		uNetINode solofstypes.NetINodeUintptr, netblockIndex int32) error
	UploadMemBlockWithSolomq(uJob solofstypes.UploadMemBlockJobUintptr,
		uploadPeerIndex int) error
}
