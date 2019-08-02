package sdfsapitypes

import "soloos/common/snettypes"

// DataNode
type GetDataNode func(peerID snettypes.PeerID) (snettypes.Peer, error)
type ChooseDataNodesForNewNetBlock func(uNetINode NetINodeUintptr) (snettypes.PeerGroup, error)
type PReadMemBlockWithDisk func(uNetINode NetINodeUintptr,
	uNetBlock NetBlockUintptr, netBlockIndex int32,
	uMemBlock MemBlockUintptr, memBlockIndex int32,
	offset uint64, length int) (int, error)
type UploadMemBlockWithDisk func(uJob UploadMemBlockJobUintptr,
	uploadPeerIndex int) error
type UploadMemBlockWithSWAL func(uJob UploadMemBlockJobUintptr,
	uploadPeerIndex int) error
