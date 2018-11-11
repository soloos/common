package snet

import (
	"soloos/snet/types"
	"soloos/util/offheap"
)

type SNetDriver struct {
	offheapDriver *offheap.OffheapDriver
	maxPeerID     int64
	peerPool      offheap.RawObjectPool
}

func (p *SNetDriver) Init(offheapDriver *offheap.OffheapDriver) error {
	var err error
	p.offheapDriver = offheapDriver
	err = p.offheapDriver.InitRawObjectPool(&p.peerPool, int(types.PeerStructSize), -1, nil, nil)
	if err != nil {
		return err
	}

	return nil
}

func (p *SNetDriver) NewPeer() types.PeerUintptr {
	var ret types.PeerUintptr
	ret = types.PeerUintptr(p.peerPool.AllocRawObject())
	return ret
}
