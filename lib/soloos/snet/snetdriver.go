package snet

import (
	"soloos/snet/types"
	"soloos/util"
	"soloos/util/offheap"
	"sync"
)

type SNetDriver struct {
	offheapDriver *offheap.OffheapDriver
	maxPeerID     int64
	peerPool      offheap.RawObjectPool
	peersRWMutex  sync.RWMutex
	peers         map[string]types.PeerUintptr
}

func (p *SNetDriver) Init(offheapDriver *offheap.OffheapDriver) error {
	var err error
	p.offheapDriver = offheapDriver
	err = p.offheapDriver.InitRawObjectPool(&p.peerPool, int(types.PeerStructSize), -1, nil, nil)
	if err != nil {
		return err
	}
	p.peers = make(map[string]types.PeerUintptr, 1024)

	return nil
}

func (p *SNetDriver) MustGetPeer(peerID *types.PeerID, addr string, protocol int) types.PeerUintptr {
	var ret types.PeerUintptr

	if peerID == nil {
		ret = types.PeerUintptr(p.peerPool.AllocRawObject())
		util.InitUUID64(&(ret.Ptr().PeerID))
		ret.Ptr().SetAddress(addr)
		ret.Ptr().ServiceProtocol = protocol
		p.peersRWMutex.Lock()
		p.peers[addr] = ret
		p.peersRWMutex.Unlock()
		return ret
	}

	p.peersRWMutex.RLock()
	ret = p.peers[addr]
	if ret != 0 {
		p.peersRWMutex.RUnlock()
		return ret
	}
	p.peersRWMutex.RUnlock()

	p.peersRWMutex.Lock()
	ret = p.peers[addr]
	if ret != 0 {
		p.peersRWMutex.Unlock()
		return ret
	}
	ret = types.PeerUintptr(p.peerPool.AllocRawObject())
	ret.Ptr().PeerID = *peerID
	ret.Ptr().SetAddress(addr)
	ret.Ptr().ServiceProtocol = types.ProtocolSRPC
	p.peers[addr] = ret
	p.peersRWMutex.Unlock()

	return ret
}
