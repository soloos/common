package snet

import (
	"soloos/common/snettypes"
	"soloos/sdbone/offheap"
)

type NetDriver struct {
	offheapDriver *offheap.OffheapDriver
	peerTable     offheap.LKVTableWithBytes64
}

func (p *NetDriver) Init(offheapDriver *offheap.OffheapDriver, name string) error {
	var err error
	p.offheapDriver = offheapDriver
	err = p.offheapDriver.InitLKVTableWithBytes64(&p.peerTable, name,
		int(snettypes.PeerStructSize), -1, offheap.DefaultKVTableSharedCount, nil)
	if err != nil {
		return err
	}

	return nil
}

func (p *NetDriver) InitPeerID(peerID *snettypes.PeerID) {
	// todo: ensure peer id unique
	snettypes.InitTmpPeerID(peerID)
}

func (p *NetDriver) GetPeer(peerID snettypes.PeerID) snettypes.PeerUintptr {
	var ret uintptr
	ret = p.peerTable.TryGetObject(peerID)
	p.peerTable.ReleaseObject(offheap.LKVTableObjectUPtrWithBytes64(ret))
	return snettypes.PeerUintptr(ret)
}

// MustGetPee return uPeer and peer is inited before
func (p *NetDriver) MustGetPeer(peerID *snettypes.PeerID, addr string, protocol int) (snettypes.PeerUintptr, bool) {
	var (
		uObject        offheap.LKVTableObjectUPtrWithBytes64
		uPeer          snettypes.PeerUintptr
		afterSetNewObj offheap.KVTableAfterSetNewObj
		newPeerID      snettypes.PeerID
		loaded         bool
	)

	if peerID == nil {
		p.InitPeerID(&newPeerID)
		peerID = &newPeerID
	}

	uObject, afterSetNewObj = p.peerTable.MustGetObject(*peerID)
	loaded = afterSetNewObj == nil
	uPeer = snettypes.PeerUintptr(uObject)
	if afterSetNewObj != nil {
		uPeer.Ptr().SetAddress(addr)
		uPeer.Ptr().ServiceProtocol = protocol
		afterSetNewObj()
	}
	p.peerTable.ReleaseObject(offheap.LKVTableObjectUPtrWithBytes64(uPeer))

	return uPeer, loaded
}
