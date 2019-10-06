package solofsapitypes

import (
	"soloos/common/snettypes"
	"soloos/common/solodbapitypes"
	"soloos/solodb/offheap"
	"unsafe"
)

const (
	NetBlockStructSize          = unsafe.Sizeof(NetBlock{})
	MaxSolodnsSizeStoreNetBlock = 8
)

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

	StorDataBackends   snettypes.PeerGroup          `db:"-"`
	IsDBMetaDataInited solodbapitypes.MetaDataState `db:"-"`

	SyncDataBackends         snettypes.TransferPeerGroup  `db:"-"`
	IsSyncDataBackendsInited solodbapitypes.MetaDataState `db:"-"`
	IsLocalDataBackendExists bool                         `db:"-"`
	IsLocalDataBackendInited solodbapitypes.MetaDataState `db:"-"`
}

func (p *NetBlock) NetINodeIDStr() string { return string(p.NetINodeID[:]) }

func (p *NetBlock) Reset() {
	p.IsDBMetaDataInited.Reset()
	p.IsSyncDataBackendsInited.Reset()
	p.IsLocalDataBackendInited.Reset()
}
