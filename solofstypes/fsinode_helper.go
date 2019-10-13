package solofstypes

// FsINode
type AllocFsINodeIno func(nsID NameSpaceID) (FsINodeIno, error)
type DeleteFsINodeByIDInDB func(nsID NameSpaceID, fsINodeIno FsINodeIno) error
type ListFsINodeByParentIDFromDB func(nsID NameSpaceID,
	parentID FsINodeIno,
	isFetchAllCols bool,
	beforeLiteralFunc func(resultCount int64) (fetchRowsLimit uint64, fetchRowsOffset uint64),
	literalFunc func(FsINodeMeta) bool,
) error
type UpdateFsINodeInDB func(nsID NameSpaceID, fsINodeMeta FsINodeMeta) error
type InsertFsINodeInDB func(nsID NameSpaceID, fsINodeMeta FsINodeMeta) error
type FetchFsINodeByIDFromDB func(nsID NameSpaceID, fsINodeIno FsINodeIno) (FsINodeMeta, error)
type FetchFsINodeByNameFromDB func(nsID NameSpaceID, parentID FsINodeIno, fsINodeName string) (FsINodeMeta, error)
