package solofstypes

// NetINode
type GetNetINode func(netINodeID NetINodeID) (NetINodeUintptr, error)
type MustGetNetINode func(netINodeID NetINodeID,
	size uint64, netBlockCap int, memBlockCap int) (NetINodeUintptr, error)
type ReleaseNetINode func(uNetINode NetINodeUintptr)
type PrepareNetINodeMetaDataOnlyLoadDB func(uNetINode NetINodeUintptr) error
type PrepareNetINodeMetaDataWithStorDB func(uNetINode NetINodeUintptr,
	size uint64, netBlockCap int, memBlockCap int) error
type NetINodeCommitSizeInDB func(uNetINode NetINodeUintptr, size uint64) error
