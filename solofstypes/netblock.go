package solofstypes

import (
	"soloos/common/snet"
	"soloos/common/solodbtypes"
	"soloos/solodb/offheap"
	"unsafe"
)

const (
	NetBlockStructSize          = unsafe.Sizeof(NetBlock{})
	MaxSolodnsSizeStoreNetBlock = 8
)

type NetBlockIndex = int32

type NetBlockUintptr uintptr

func (u NetBlockUintptr) Ptr() *NetBlock { return (*NetBlock)(unsafe.Pointer(u)) }

type NetBlockMeta struct {
	NetINodeID      NetINodeID `db:"netinode_id"`
	IndexInNetINode int32      `db:"index_in_netinode"`
	Len             int        `db:"netblock_len"`
	Cap             int        `db:"netblock_cap"`
}

type NetBlock struct {
	offheap.LKVTableObjectWithBytes68 `db:"-"`

	NetBlockMeta

	StorDataBackends   snet.PeerGroup            `db:"-"`
	IsDBMetaDataInited solodbtypes.MetaDataState `db:"-"`

	SyncDataBackends         snet.TransferPeerGroup    `db:"-"`
	IsSyncDataBackendsInited solodbtypes.MetaDataState `db:"-"`
	IsLocalDataBackendExists bool                      `db:"-"`
	IsLocalDataBackendInited solodbtypes.MetaDataState `db:"-"`
}

func (p *NetBlock) NetINodeIDStr() string { return string(p.NetINodeID[:]) }

func (p *NetBlock) Reset() {
	p.IsDBMetaDataInited.Reset()
	p.IsSyncDataBackendsInited.Reset()
	p.IsLocalDataBackendInited.Reset()
}
