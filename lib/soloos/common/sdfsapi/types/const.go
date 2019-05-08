package types

import (
	snettypes "soloos/common/snet/types"
	"unsafe"
)

const (
	DefaultSDFSRPCNetwork  = "tcp"
	DefaultSDFSRPCProtocol = snettypes.ProtocolSRPC

	UUintptrSize = unsafe.Sizeof(uintptr(0))
	UintptrSize  = int(unsafe.Sizeof(uintptr(0)))
	Int32Size    = int(unsafe.Sizeof(int32(0)))

	FS_MAX_PATH_LENGTH = 2048
	FS_MAX_NAME_LENGTH = 128

	MaxFsINodeNameLen = FS_MAX_NAME_LENGTH
)