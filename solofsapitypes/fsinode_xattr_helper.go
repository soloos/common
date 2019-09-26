package solofsapitypes

// FsINodeXAttr
type DeleteFIXAttrInDB func(nameSpaceID NameSpaceID, fsINodeID FsINodeID) error
type ReplaceFIXAttrInDB func(nameSpaceID NameSpaceID, fsINodeID FsINodeID, xattr FsINodeXAttr) error
type GetFIXAttrByInoFromDB func(nameSpaceID NameSpaceID, fsINodeID FsINodeID) (FsINodeXAttr, error)
