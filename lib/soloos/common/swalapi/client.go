package swalapi

import (
	sdfsapitypes "soloos/common/sdfsapi/types"
)

type SendLogSig = bool

type Client interface {
	SetSDFSClient(itSDFSClient interface{}) error
	UploadMemBlockWithSWAL(uJob sdfsapitypes.UploadMemBlockJobUintptr,
		uploadPeerIndex int, transferPeersCount int) error
}
