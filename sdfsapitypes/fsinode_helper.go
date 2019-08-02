package sdfsapitypes

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
