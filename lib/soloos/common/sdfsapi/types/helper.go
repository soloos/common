package types

import (
	snettypes "soloos/common/snet/types"
)

// DataNode
type GetDataNode func(peerID snettypes.PeerID) snettypes.PeerUintptr
type ChooseDataNodesForNewNetBlock func(uNetINode NetINodeUintptr) (snettypes.PeerGroup, error)

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
type DeleteFsINodeByIDInDB func(fsINodeID FsINodeID) error
type ListFsINodeByParentIDFromDB func(parentID FsINodeID,
	isFetchAllCols bool,
	beforeLiteralFunc func(resultCount int) (fetchRowsLimit uint64, fetchRowsOffset uint64),
	literalFunc func(FsINodeMeta) bool,
) error
type UpdateFsINodeInDB func(fsINodeMeta *FsINodeMeta) error
type InsertFsINodeInDB func(fsINodeMeta *FsINodeMeta) error
type FetchFsINodeByIDFromDB func(pFsINodeMeta *FsINodeMeta) error
type FetchFsINodeByNameFromDB func(pFsINodeMeta *FsINodeMeta) error

// FsINodeXAttr
type DeleteFIXAttrInDB func(fsINodeID FsINodeID) error
type ReplaceFIXAttrInDB func(fsINodeID FsINodeID, xattr FsINodeXAttr) error
type GetFIXAttrByInoFromDB func(fsINodeID FsINodeID) (FsINodeXAttr, error)
