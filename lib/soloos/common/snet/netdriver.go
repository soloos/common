package snet

import (
	"soloos/common/snet/types"
	"soloos/common/util"
	"soloos/common/util/offheap"
	"sync"
)

type NetDriver struct {
	offheapDriver *offheap.OffheapDriver
	maxPeerID     int64
	peerPool      offheap.RawObjectPool
	peersRWMutex  sync.RWMutex
	peers         map[types.PeerID]types.PeerUintptr
}

func (p *NetDriver) Init(offheapDriver *offheap.OffheapDriver) error {
	var err error
	p.offheapDriver = offheapDriver
	err = p.offheapDriver.InitRawObjectPool(&p.peerPool, int(types.PeerStructSize), -1, nil, nil)
	if err != nil {
		return err
	}
	p.peers = make(map[types.PeerID]types.PeerUintptr, 1024)

	return nil
}

func (p *NetDriver) GetPeer(peerID types.PeerID) types.PeerUintptr {
	var ret types.PeerUintptr
	p.peersRWMutex.RLock()
	ret = p.peers[peerID]
	p.peersRWMutex.RUnlock()
	return ret
}

// MustGetPee return uPeer and peer is inited before
func (p *NetDriver) MustGetPeer(peerID *types.PeerID, addr string, protocol int) (types.PeerUintptr, bool) {
	var ret types.PeerUintptr

	if peerID == nil {
		ret = types.PeerUintptr(p.peerPool.AllocRawObject())
		util.InitUUID64(&(ret.Ptr().PeerID))
		ret.Ptr().SetAddress(addr)
		ret.Ptr().ServiceProtocol = protocol
		p.peersRWMutex.Lock()
		p.peers[ret.Ptr().PeerID] = ret
		p.peersRWMutex.Unlock()
		return ret, false
	}

	p.peersRWMutex.RLock()
	ret = p.peers[*peerID]
	if ret != 0 {
		p.peersRWMutex.RUnlock()
		return ret, true
	}
	p.peersRWMutex.RUnlock()

	p.peersRWMutex.Lock()
	ret = p.peers[*peerID]
	if ret != 0 {
		p.peersRWMutex.Unlock()
		return ret, true
	}
	ret = types.PeerUintptr(p.peerPool.AllocRawObject())
	ret.Ptr().PeerID = *peerID
	ret.Ptr().SetAddress(addr)
	ret.Ptr().ServiceProtocol = protocol
	p.peers[*peerID] = ret
	p.peersRWMutex.Unlock()

	return ret, false
}
