package solofsapitypes

// FsINodeXAttr
type DeleteFIXAttrInDB func(nsID NameSpaceID, fsINodeID FsINodeID) error
type ReplaceFIXAttrInDB func(nsID NameSpaceID, fsINodeID FsINodeID, xattr FsINodeXAttr) error
type GetFIXAttrByInoFromDB func(nsID NameSpaceID, fsINodeID FsINodeID) (FsINodeXAttr, error)
