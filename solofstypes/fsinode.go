package solofstypes

import (
	"soloos/common/solodbtypes"
	"soloos/solodb/offheap"
	"unsafe"
)

type FsINodeIno = uint64
type DirTreeTime = uint64
type DirTreeTimeNsec = uint32

const (
	MaxUint64             = ^uint64(0)
	ZombieFsINodeParentID = FsINodeIno(0)
	RootFsINodeParentID   = FsINodeIno(0)
	RootFsINodeIno         = FsINodeIno(1)
	FsINodeStructSize     = unsafe.Sizeof(FsINode{})
	MaxFsINodeIno          = MaxUint64
)

type FsINodeUintptr uintptr

func (u FsINodeUintptr) Ptr() *FsINode { return (*FsINode)(unsafe.Pointer(u)) }

type FsINodeMeta struct {
	LastModifyACMTime int64
	LoadInMemAt       int64

	NameSpaceID  NameSpaceID
	Ino          FsINodeIno
	HardLinkIno  FsINodeIno
	NetINodeID   NetINodeID
	ParentID     FsINodeIno
	NameBytesLen int
	NameBytes    [MaxFsINodeNameLen]byte
	Type         int
	Atime        DirTreeTime
	Ctime        DirTreeTime
	Mtime        DirTreeTime
	Atimensec    DirTreeTimeNsec
	Ctimensec    DirTreeTimeNsec
	Mtimensec    DirTreeTimeNsec
	Mode         uint32
	Nlink        int32
	Uid          uint32
	Gid          uint32
	Rdev         uint32
}

func (p *FsINodeMeta) SetName(nameStr string) {
	p.NameBytesLen = len(nameStr)
	if p.NameBytesLen > MaxFsINodeNameLen {
		p.NameBytesLen = MaxFsINodeNameLen
	}
	copy(p.NameBytes[:p.NameBytesLen], []byte(nameStr))
}

func (p *FsINodeMeta) Name() string {
	return string(p.NameBytes[:p.NameBytesLen])
}

func (p *FsINodeMeta) NameLen() int {
	return p.NameBytesLen
}

type FsINode struct {
	offheap.LKVTableObjectWithUint64 `db:"-"`
	Meta                             FsINodeMeta
	UNetINode                        NetINodeUintptr
	IsDBMetaDataInited               solodbtypes.MetaDataState `db:"-"`
}

func (p *FsINode) Reset() {
	p.IsDBMetaDataInited.Reset()
}
