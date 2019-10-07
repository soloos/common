package solofsapitypes

import "soloos/common/snet"

// Solodn
type GetSolodn func(peerID snet.PeerID) (snet.Peer, error)
type ChooseSolodnsForNewNetBlock func(uNetINode NetINodeUintptr) (snet.PeerGroup, error)
type PReadMemBlockWithDisk func(uNetINode NetINodeUintptr,
	uNetBlock NetBlockUintptr, netBlockIndex int32,
	uMemBlock MemBlockUintptr, memBlockIndex int32,
	offset uint64, length int) (int, error)
type UploadMemBlockWithDisk func(uJob UploadMemBlockJobUintptr,
	uploadPeerIndex int) error
type UploadMemBlockWithSolomq func(uJob UploadMemBlockJobUintptr,
	uploadPeerIndex int) error
