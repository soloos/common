package sdfsapitypes

import (
	"soloos/common/snettypes"
)

// DataNode
type GetDataNode func(peerID snettypes.PeerID) (snettypes.Peer, error)
type ChooseDataNodesForNewNetBlock func(uNetINode NetINodeUintptr) (snettypes.PeerGroup, error)
type PReadMemBlockWithDisk func(uNetINode NetINodeUintptr,
	uNetBlock NetBlockUintptr, netBlockIndex int32,
	uMemBlock MemBlockUintptr, memBlockIndex int32,
	offset uint64, length int) (int, error)
type UploadMemBlockWithDisk func(uJob UploadMemBlockJobUintptr,
	uploadPeerIndex int, transferPeersCount int) error
type UploadMemBlockWithSWAL func(uJob UploadMemBlockJobUintptr,
	uploadPeerIndex int, transferPeersCount int) error

// NetINode
type GetNetINode func(netINodeID NetINodeID) (NetINodeUintptr, error)
type MustGetNetINode func(netINodeID NetINodeID,
	size uint64, netBlockCap int, memBlockCap int) (NetINodeUintptr, error)
type ReleaseNetINode func(uNetINode NetINodeUintptr)
type PrepareNetINodeMetaDataOnlyLoadDB func(uNetINode NetINodeUintptr) error
type PrepareNetINodeMetaDataWithStorDB func(uNetINode NetINodeUintptr,
	size uint64, netBlockCap int, memBlockCap int) error
type NetINodeCommitSizeInDB func(uNetINode NetINodeUintptr, size uint64) error

// FsINode
type AllocFsINodeID func() FsINodeID
type DeleteFsINodeByIDInDB func(nameSpaceID NameSpaceID, fsINodeID FsINodeID) error
type ListFsINodeByParentIDFromDB func(nameSpaceID NameSpaceID,
	parentID FsINodeID,
	isFetchAllCols bool,
	beforeLiteralFunc func(resultCount int) (fetchRowsLimit uint64, fetchRowsOffset uint64),
	literalFunc func(FsINodeMeta) bool,
) error
type UpdateFsINodeInDB func(nameSpaceID NameSpaceID, fsINodeMeta FsINodeMeta) error
type InsertFsINodeInDB func(nameSpaceID NameSpaceID, fsINodeMeta FsINodeMeta) error
type FetchFsINodeByIDFromDB func(nameSpaceID NameSpaceID, fsINodeID FsINodeID) (FsINodeMeta, error)
type FetchFsINodeByNameFromDB func(nameSpaceID NameSpaceID, parentID FsINodeID, fsINodeName string) (FsINodeMeta, error)

// FsINodeXAttr
type DeleteFIXAttrInDB func(nameSpaceID NameSpaceID, fsINodeID FsINodeID) error
type ReplaceFIXAttrInDB func(nameSpaceID NameSpaceID, fsINodeID FsINodeID, xattr FsINodeXAttr) error
type GetFIXAttrByInoFromDB func(nameSpaceID NameSpaceID, fsINodeID FsINodeID) (FsINodeXAttr, error)
