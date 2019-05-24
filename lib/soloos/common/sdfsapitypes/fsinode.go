package sdfsapitypes

import (
	"soloos/common/sdbapitypes"
	"soloos/sdbone/offheap"
	"unsafe"
)

type FsINodeID = uint64
type DirTreeTime = uint64
type DirTreeTimeNsec = uint32

const (
	MaxUint64             = ^uint64(0)
	ZombieFsINodeParentID = FsINodeID(0)
	RootFsINodeParentID   = FsINodeID(0)
	RootFsINodeID         = FsINodeID(1)
	FsINodeStructSize     = unsafe.Sizeof(FsINode{})
	MaxFsINodeID          = MaxUint64
)

type FsINodeUintptr uintptr

func (u FsINodeUintptr) Ptr() *FsINode { return (*FsINode)(unsafe.Pointer(u)) }

type FsINodeMeta struct {
	LastModifyACMTime int64
	LoadInMemAt       int64

	Ino         FsINodeID
	HardLinkIno FsINodeID
	NetINodeID  NetINodeID
	ParentID    FsINodeID
	nameLen     int
	nameBytes   [MaxFsINodeNameLen]byte
	Type        int
	Atime       DirTreeTime
	Ctime       DirTreeTime
	Mtime       DirTreeTime
	Atimensec   DirTreeTimeNsec
	Ctimensec   DirTreeTimeNsec
	Mtimensec   DirTreeTimeNsec
	Mode        uint32
	Nlink       int32
	Uid         uint32
	Gid         uint32
	Rdev        uint32
}

func (p *FsINodeMeta) SetName(nameStr string) {
	p.nameLen = len(nameStr)
	if p.nameLen > MaxFsINodeNameLen {
		p.nameLen = MaxFsINodeNameLen
	}
	copy(p.nameBytes[:p.nameLen], []byte(nameStr))
}

func (p *FsINodeMeta) Name() string {
	return string(p.nameBytes[:p.nameLen])
}

func (p *FsINodeMeta) NameLen() int {
	return p.nameLen
}

type FsINode struct {
	offheap.LKVTableObjectWithUint64 `db:"-"`
	Meta                             FsINodeMeta
	UNetINode                        NetINodeUintptr
	IsDBMetaDataInited               sdbapitypes.MetaDataState `db:"-"`
}

func (p *FsINode) Reset() {
	p.IsDBMetaDataInited.Reset()
}
