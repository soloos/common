package solofstypes

// FsINodeXAttr
type DeleteFIXAttrInDB func(nsID NameSpaceID, fsINodeIno FsINodeIno) error
type ReplaceFIXAttrInDB func(nsID NameSpaceID, fsINodeIno FsINodeIno, xattr FsINodeXAttr) error
type GetFIXAttrByInoFromDB func(nsID NameSpaceID, fsINodeIno FsINodeIno) (FsINodeXAttr, error)
