package solofsapitypes

// FsINode
type AllocFsINodeID func(nsID NameSpaceID) (FsINodeID, error)
type DeleteFsINodeByIDInDB func(nsID NameSpaceID, fsINodeID FsINodeID) error
type ListFsINodeByParentIDFromDB func(nsID NameSpaceID,
	parentID FsINodeID,
	isFetchAllCols bool,
	beforeLiteralFunc func(resultCount int64) (fetchRowsLimit uint64, fetchRowsOffset uint64),
	literalFunc func(FsINodeMeta) bool,
) error
type UpdateFsINodeInDB func(nsID NameSpaceID, fsINodeMeta FsINodeMeta) error
type InsertFsINodeInDB func(nsID NameSpaceID, fsINodeMeta FsINodeMeta) error
type FetchFsINodeByIDFromDB func(nsID NameSpaceID, fsINodeID FsINodeID) (FsINodeMeta, error)
type FetchFsINodeByNameFromDB func(nsID NameSpaceID, parentID FsINodeID, fsINodeName string) (FsINodeMeta, error)
