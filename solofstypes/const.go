package solofstypes

import (
	"soloos/common/fsapi"
	"soloos/common/snet"
	"syscall"
	"unsafe"
)

const (
	DefaultSolofsRPCNetwork = "tcp"

	UUintptrSize = unsafe.Sizeof(uintptr(0))
	UintptrSize  = int(unsafe.Sizeof(uintptr(0)))
	Int32Size    = int(unsafe.Sizeof(int32(0)))

	FS_MAX_PATH_LENGTH = 2048
	FS_MAX_NAME_LENGTH = 128

	MaxFsINodeNameLen = FS_MAX_NAME_LENGTH
)

var (
	DefaultSolofsRPCProtocol = snet.ProtocolSolofs
)

const (
	FSINODE_TYPE_FILE = iota
	FSINODE_TYPE_DIR
	FSINODE_TYPE_HARD_LINK
	FSINODE_TYPE_SOFT_LINK
	FSINODE_TYPE_FIFO
	FSINODE_TYPE_UNKOWN
)

const (
	FS_RDEV = 0

	FS_EEXIST       = fsapi.Status(syscall.EEXIST)
	FS_ENOTEMPTY    = fsapi.Status(syscall.ENOTEMPTY)
	FS_ENAMETOOLONG = fsapi.Status(syscall.ENAMETOOLONG)

	FS_INODE_LOCK_SH = syscall.LOCK_SH
	FS_INODE_LOCK_EX = syscall.LOCK_EX
	FS_INODE_LOCK_UN = syscall.LOCK_UN

	FS_XATTR_SOFT_LNKMETA_KEY = "solofs.soft.link.metadata"
)

const (
	FS_PERM_SETUID uint32 = 1 << (12 - 1 - iota)
	FS_PERM_SETGID
	FS_PERM_STICKY
	FS_PERM_USER_READ
	FS_PERM_USER_WRITE
	FS_PERM_USER_EXECUTE
	FS_PERM_GROUP_READ
	FS_PERM_GROUP_WRITE
	FS_PERM_GROUP_EXECUTE
	FS_PERM_OTHER_READ
	FS_PERM_OTHER_WRITE
	FS_PERM_OTHER_EXECUTE
)
