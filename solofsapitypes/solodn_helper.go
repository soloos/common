package solofsapitypes

import "soloos/common/snettypes"

// Solodn
type GetSolodn func(peerID snettypes.PeerID) (snettypes.Peer, error)
type ChooseSolodnsForNewNetBlock func(uNetINode NetINodeUintptr) (snettypes.PeerGroup, error)
type PReadMemBlockWithDisk func(uNetINode NetINodeUintptr,
	uNetBlock NetBlockUintptr, netBlockIndex int32,
	uMemBlock MemBlockUintptr, memBlockIndex int32,
	offset uint64, length int) (int, error)
type UploadMemBlockWithDisk func(uJob UploadMemBlockJobUintptr,
	uploadPeerIndex int) error
type UploadMemBlockWithSOLOMQ func(uJob UploadMemBlockJobUintptr,
	uploadPeerIndex int) error
