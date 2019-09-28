package solofsapitypes

import (
	"soloos/common/solodbapitypes"
	"soloos/common/util"
	"soloos/solodb/offheap"
	"sync"
	"unsafe"
)

const (
	NetINodeIDBytesNum = 64
	NetINodeIDSize     = int(unsafe.Sizeof([NetINodeIDBytesNum]byte{}))
	NetINodeStructSize = unsafe.Sizeof(NetINode{})
)

var (
	ZeroNetINodeID NetINodeID
)

type NetINodeID [NetINodeIDBytesNum]byte
type NetINodeUintptr uintptr

func InitTmpNetINodeID(netINodeID *NetINodeID) {
	util.InitUUID64((*[64]byte)(netINodeID))
}

func (p NetINodeID) Str() string {
	return string(p[:])
}

func init() {
	copy(ZeroNetINodeID[:], ([]byte("0000000000000000000000000000000000000000000000000000000000000000")[:64]))
}

func (u NetINodeUintptr) Ptr() *NetINode { return (*NetINode)(unsafe.Pointer(u)) }

type NetINodeMeta struct {
	NetINodeID  NetINodeID `db:"netinode_id"`
	Size        uint64     `db:"netinode_size"`
	NetBlockCap int        `db:"netblock_cap"`
	MemBlockCap int        `db:"memblock_cap"`
}

type NetINode struct {
	offheap.LKVTableObjectWithBytes64 `db:"-"`

	MemBlockPlacementPolicy MemBlockPlacementPolicy `db:"-"`

	LastCommitSize uint64 `db:"-"`

	NetINodeMeta

	WriteDataRWMutex   sync.RWMutex                 `db:"-"`
	SyncDataSig        util.RawWaitGroup            `db:"-"`
	LastSyncDataError  error                        `db:"-"`
	IsDBMetaDataInited solodbapitypes.MetaDataState `db:"-"`
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
