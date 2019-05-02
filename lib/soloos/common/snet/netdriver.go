package snet

import (
	"soloos/common/snet/types"
	"soloos/common/util"
	"soloos/sdbone/offheap"
)

type NetDriver struct {
	offheapDriver *offheap.OffheapDriver
	peers         offheap.LKVTableWithBytes64
}

func (p *NetDriver) Init(offheapDriver *offheap.OffheapDriver, name string) error {
	var err error
	p.offheapDriver = offheapDriver
	err = p.offheapDriver.InitLKVTableWithBytes64(&p.peers, name,
		int(types.PeerStructSize), -1, offheap.DefaultKVTableSharedCount, nil)
	if err != nil {
		return err
	}

	return nil
}

func (p *NetDriver) InitPeerID(peerID *types.PeerID) {
	// todo: ensure peer id unique
	util.InitUUID64(peerID)
}

func (p *NetDriver) GetPeer(peerID types.PeerID) types.PeerUintptr {
	var ret uintptr
	ret = p.peers.TryGetObjectWithAcquire(peerID)
	if ret != 0 {
		p.peers.ReleaseObject(offheap.LKVTableObjectUPtrWithBytes64(ret))
	}
	return types.PeerUintptr(ret)
}

func (p *NetDriver) AllocPeer(addr string, protocol int) types.PeerUintptr {
	var (
		peerID types.PeerID
		uPeer  types.PeerUintptr
	)
	p.InitPeerID(&peerID)
	uPeer, _ = p.RegisterPeer(&peerID, addr, protocol)
	return uPeer
}

func (p *NetDriver) RegisterPeer(peerID *types.PeerID, addr string, protocol int) (types.PeerUintptr, bool) {
	return p.MustGetPeer(peerID, addr, protocol)
}

// MustGetPee return uPeer and peer is inited before
func (p *NetDriver) MustGetPeer(peerID *types.PeerID, addr string, protocol int) (types.PeerUintptr, bool) {
	var (
		u              uintptr
		uPeer          types.PeerUintptr
		afterSetNewObj offheap.KVTableAfterSetNewObj
		loaded         bool
	)

	u, afterSetNewObj = p.peers.MustGetObjectWithAcquire(*peerID)
	loaded = afterSetNewObj == nil
	uPeer = types.PeerUintptr(u)
	if afterSetNewObj != nil {
		uPeer.Ptr().SetAddress(addr)
		uPeer.Ptr().ServiceProtocol = protocol
		afterSetNewObj()
	}
	p.peers.ReleaseObject(offheap.LKVTableObjectUPtrWithBytes64(uPeer))

	return uPeer, loaded
}
