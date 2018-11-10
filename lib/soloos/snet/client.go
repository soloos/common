package snet

import (
	"soloos/snet/types"
	"soloos/util/offheap"
)

type Client struct {
	offheapDriver *offheap.OffheapDriver
	peerPool      offheap.RawObjectPool
	netDrivers    map[string]NetDriver
}

func (p *Client) RawChunkPoolInvokeReleaseRawChunk() {
	panic("not support")
}

func (p *Client) RawChunkPoolInvokePrepareNewRawChunk(uRawChunk uintptr) {
}

func (p *Client) Init(offheapDriver *offheap.OffheapDriver) error {
	var err error
	p.offheapDriver = offheapDriver
	err = p.offheapDriver.InitRawObjectPool(&p.peerPool,
		int(types.PeerStructSize), -1,
		p.RawChunkPoolInvokePrepareNewRawChunk, p.RawChunkPoolInvokeReleaseRawChunk)
	if err != nil {
		return err
	}

	return nil
}

func (p *Client) RegisterPeer(peer types.Peer) types.PeerUintptr {
	u, _ := p.peerPool.MustGetRawObject(peer.ID)
	uPeer := (types.PeerUintptr)(u)
	*uPeer.Ptr() = peer
	return uPeer
}

func (p *Client) Write(uPeer types.PeerUintptr) error {
	return nil
}
