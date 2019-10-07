package fsapi

import (
	. "soloos/common/fsapitypes"
	"soloos/common/snet"
	"soloos/common/solofsapitypes"
)

// PosixFs is an interface close to the FUSE wire solofsprotocol.
//
// Unless you really know what you are doing, you should not implement
// this, but rather the nodefs.Node or pathfs.FileSystem interfaces; the
// details of getting interactions with open files, renames, and threading
// right etc. are somewhat tricky and not very interesting.
//
// Each FUSE request results in a corresponding method called by Server.
// Several calls may be made simultaneously, because the server typically calls
// each method in separate goroutine.
//
// A null implementation is provided by NewDefaultPosixFs.
type PosixFs interface {
	String() string

	// If called, provide debug output through the log package.
	SetDebug(debug bool)

	// Lookup is called by the kernel when the VFS wants to know
	// about a file inside a directory. Many lookup calls can
	// occur in parallel, but only one call happens for each (dir,
	// name) pair.
	Lookup(header *InHeader, name string, out *EntryOut) (status Status)

	// Forget is called when the kernel discards entries from its
	// dentry cache. This happens on unmount, and when the kernel
	// is short on memory. Since it is not guaranteed to occur at
	// any moment, and since there is no return value, Forget
	// should not do I/O, as there is no channel to report back
	// I/O errors.
	Forget(nodeid, nlookup uint64)

	// Attributes.
	GetAttr(input *GetAttrIn, out *AttrOut) (code Status)
	SetAttr(input *SetAttrIn, out *AttrOut) (code Status)

	// Modifying structure.
	Mknod(input *MknodIn, name string, out *EntryOut) (code Status)
	Mkdir(input *MkdirIn, name string, out *EntryOut) (code Status)
	Unlink(header *InHeader, name string) (code Status)
	Rmdir(header *InHeader, name string) (code Status)
	Rename(input *RenameIn, oldName string, newName string) (code Status)
	Link(input *LinkIn, filename string, out *EntryOut) (code Status)

	Symlink(header *InHeader, pointedTo string, linkName string, out *EntryOut) (code Status)
	Readlink(header *InHeader) (out []byte, code Status)
	Access(input *AccessIn) (code Status)

	// Extended attributes.
	GetXAttrSize(header *InHeader, attr string) (sz int, code Status)
	GetXAttrData(header *InHeader, attr string) (data []byte, code Status)
	ListXAttr(header *InHeader) (attributes []byte, code Status)
	SetXAttr(input *SetXAttrIn, attr string, data []byte) Status
	RemoveXAttr(header *InHeader, attr string) (code Status)

	// File handling.
	Create(input *CreateIn, name string, out *CreateOut) (code Status)
	Open(input *OpenIn, out *OpenOut) (status Status)
	Read(input *ReadIn, buf []byte) (ReadResult, Status)

	// File locking
	GetLk(input *LkIn, out *LkOut) (code Status)
	SetLk(input *LkIn) (code Status)
	SetLkw(input *LkIn) (code Status)

	Release(input *ReleaseIn)
	Write(input *WriteIn, data []byte) (written uint32, code Status)
	Flush(input *FlushIn) Status
	Fsync(input *FsyncIn) (code Status)
	Fallocate(input *FallocateIn) (code Status)

	// Directory handling
	OpenDir(input *OpenIn, out *OpenOut) (status Status)
	ReadDir(input *ReadIn, out *DirEntryList) Status
	ReadDirPlus(input *ReadIn, out *DirEntryList) Status
	ReleaseDir(input *ReleaseIn)
	FsyncDir(input *FsyncIn) (code Status)

	//
	StatFs(input *InHeader, out *StatfsOut) (code Status)

	// This is called on processing the first request. The
	// filesystem implementation can use the server argument to
	// talk back to the kernel (through notify methods).
	FsInit()

	// other
	FetchFsINodeByID(pFsINodeMeta *solofsapitypes.FsINodeMeta, fsINodeID solofsapitypes.FsINodeID) error
	FetchFsINodeByPath(pFsINodeMeta *solofsapitypes.FsINodeMeta, fsINodePath string) error
	ListFsINodeByParentPath(parentPath string,
		isFetchAllCols bool,
		beforeLiteralFunc func(resultCount int64) (fetchRowsLimit uint64, fetchRowsOffset uint64),
		literalFunc func(solofsapitypes.FsINodeMeta) bool,
	) error
	DeleteFsINodeByPath(fsINodePath string) error
	RenameWithFullPath(oldFsINodeName, newFsINodePath string) error

	FdTableAllocFd(fsINodeID solofsapitypes.FsINodeID) solofsapitypes.FsINodeFileHandlerID
	FdTableGetFd(fdID solofsapitypes.FsINodeFileHandlerID) solofsapitypes.FsINodeFileHandler
	FdTableFdAddAppendPosition(fdID solofsapitypes.FsINodeFileHandlerID, delta uint64)
	FdTableFdAddReadPosition(fdID solofsapitypes.FsINodeFileHandlerID, delta uint64)

	SimpleOpenFile(fsINodePath string, netBlockCap int, memBlockCap int) (solofsapitypes.FsINodeMeta, error)
	SimpleWriteWithMem(fsINodeID solofsapitypes.FsINodeID, data []byte, offset uint64) error
	SimpleReadWithMem(fsINodeID solofsapitypes.FsINodeID, data []byte, offset uint64) (int, error)
	SimpleFlush(fsINodeID solofsapitypes.FsINodeID) error

	SimpleMkdirAll(perms uint32, fsINodePath string, uid uint32, gid uint32) Status
	SimpleMkdir(fsINodeMeta *solofsapitypes.FsINodeMeta,
		fsINodeID *solofsapitypes.FsINodeID, parentID solofsapitypes.FsINodeID,
		perms uint32, name string,
		uid uint32, gid uint32, rdev uint32) Status

	SetNetINodeBlockPlacement(netINodeID solofsapitypes.NetINodeID, policy solofsapitypes.MemBlockPlacementPolicy) error

	GetFsINodeByID(fsINodeID solofsapitypes.FsINodeID) (solofsapitypes.FsINodeUintptr, error)
	ReleaseFsINode(uFsINode solofsapitypes.FsINodeUintptr)

	GetNetINode(netINodeID solofsapitypes.NetINodeID) (solofsapitypes.NetINodeUintptr, error)
	ReleaseNetINode(uNetINode solofsapitypes.NetINodeUintptr)
	NetINodePWriteWithNetQuery(uNetINode solofsapitypes.NetINodeUintptr,
		netQuery *snet.NetQuery, dataLength int, offset uint64) error
	NetINodePWriteWithMem(uNetINode solofsapitypes.NetINodeUintptr,
		data []byte, offset uint64) error

	MustGetNetBlock(uNetINode solofsapitypes.NetINodeUintptr, netBlockIndex int32) (solofsapitypes.NetBlockUintptr, error)
	ReleaseNetBlock(uNetBlock solofsapitypes.NetBlockUintptr)
	NetBlockSetPReadMemBlockWithDisk(preadWithDisk solofsapitypes.PReadMemBlockWithDisk)
	NetBlockSetUploadMemBlockWithDisk(uploadMemBlockWithDisk solofsapitypes.UploadMemBlockWithDisk)

	MustGetMemBlockWithReadAcquire(uNetINode solofsapitypes.NetINodeUintptr,
		memBlockIndex int32) (solofsapitypes.MemBlockUintptr, bool)
	ReleaseMemBlockWithReadRelease(uMemBlock solofsapitypes.MemBlockUintptr)
}
