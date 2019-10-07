package solofsapitypes

import (
	"soloos/common/snet"
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
