package types

import (
	"soloos/sdbone/offheap"
	"sync"
	"unsafe"
)

const (
	NetINodeIDBytesNum = 64
	NetINodeIDSize     = int(unsafe.Sizeof([NetINodeIDBytesNum]byte{}))
	NetINodeStructSize = unsafe.Sizeof(NetINode{})
)

type NetINodeID = [NetINodeIDBytesNum]byte
type NetINodeUintptr uintptr

var (
	ZeroNetINodeID NetINodeID
)

func init() {
	copy(ZeroNetINodeID[:], ([]byte("0000000000000000000000000000000000000000000000000000000000000000")[:64]))
}

func (u NetINodeUintptr) Ptr() *NetINode { return (*NetINode)(unsafe.Pointer(u)) }

type NetINode struct {
	offheap.LKVTableObjectWithBytes64 `db:"-"`

	LastCommitSize uint64 `db:"-"`

	ID                  NetINodeID     `db:"netinode_id"`
	Size                uint64         `db:"netinode_size"`
	NetBlockCap         int            `db:"netblock_cap"`
	MemBlockCap         int            `db:"memblock_cap"`
	WriteDataRWMutex    sync.RWMutex   `db:"-"`
	SyncDataSig         sync.WaitGroup `db:"-"`
	LastSyncDataError   error          `db:"-"`
	DBMetaDataInitMutex sync.Mutex     `db:"-"`
	IsDBMetaDataInited  MetaDataState  `db:"-"`
}

func (p *NetINode) IDStr() string { return string(p.ID[:]) }

func (p *NetINode) Reset() {
	p.IsDBMetaDataInited.Reset()
}

// TODO return real blocks
func (p *NetINode) GetBlocks() uint64 {
	if p.MemBlockCap == 0 {
		return 0
	}
	return p.Size / uint64(p.NetBlockCap)
}
