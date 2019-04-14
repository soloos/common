package types

import (
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

type FsINode struct {
	HSharedPointer    offheap.HSharedPointer
	LastModifyACMTime int64
	LoadInMemAt       int64

	Ino         FsINodeID
	HardLinkIno FsINodeID
	NetINodeID  NetINodeID
	ParentID    FsINodeID
	Name        string
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
	UNetINode   NetINodeUintptr
}

func (p *FsINode) Reset() {
	p.HSharedPointer.Reset()
}
