package snet

import (
	"soloos/common/snettypes"
	"soloos/sdbone/offheap"
	"sync"
)

type SRPCClientDriver struct {
	offheapDriver      *offheap.OffheapDriver
	netDriver          *NetDriver
	netConnReadSigPool offheap.RawObjectPool
	clientRWMutex      sync.RWMutex
	clients            map[snettypes.PeerID]*SRPCClient
}

func (p *SRPCClientDriver) Init(offheapDriver *offheap.OffheapDriver, netDriver *NetDriver) error {
	var err error

	p.offheapDriver = offheapDriver
	p.netDriver = netDriver
	err = p.netConnReadSigPool.Init(int(offheap.MutexStrutSize), -1, nil, nil)
	if err != nil {
		return err
	}

	p.clients = make(map[snettypes.PeerID]*SRPCClient)

	return nil
}

func (p *SRPCClientDriver) getClient(peerID snettypes.PeerID) (*SRPCClient, error) {
	var (
		ret  *SRPCClient
		peer snettypes.Peer
		err  error
	)

	p.clientRWMutex.RLock()
	ret = p.clients[peerID]
	p.clientRWMutex.RUnlock()

	if ret == nil {
		ret = p.clients[peerID]
		if ret == nil {
			peer, err = p.netDriver.GetPeer(peerID)
			if err == nil {
				ret, err = p.registerClient(peer)
			}
		}
	}

	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (p *SRPCClientDriver) registerClient(peer snettypes.Peer) (*SRPCClient, error) {
	var (
		client *SRPCClient
		err    error
	)

	p.clientRWMutex.Lock()
	client = p.clients[peer.ID]
	if client != nil {
		goto GET_CLIENT_DONE
	}

	client = &SRPCClient{}
	err = client.Init(p, peer.AddressStr())
	if err != nil {
		client = nil
		goto GET_CLIENT_DONE
	}

	err = client.Start()
	if err != nil {
		client = nil
		goto GET_CLIENT_DONE
	}

	p.clients[peer.ID] = client

GET_CLIENT_DONE:
	p.clientRWMutex.Unlock()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (p *SRPCClientDriver) CloseClient(peerID snettypes.PeerID) error {
	var (
		client = p.clients[peerID]
		err    error
	)

	err = p.sendCloseCmd(client)
	if err != nil {
		return err
	}

	err = client.Close(snettypes.ErrClosedByUser)
	if err != nil {
		return err
	}

	return nil
}
