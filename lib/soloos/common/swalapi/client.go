package swalapi

import (
	"soloos/common/sdfsapitypes"
)

type SendLogSig = bool

type Client interface {
	SetSDFSClient(itSDFSClient interface{}) error
	PrepareNetBlockMetaData(uNetBlock sdfsapitypes.NetBlockUintptr,
		uNetINode sdfsapitypes.NetINodeUintptr, netblockIndex int32) error
	UploadMemBlockWithSWAL(uJob sdfsapitypes.UploadMemBlockJobUintptr,
		uploadPeerIndex int, transferPeersCount int) error
}
