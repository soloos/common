package snet

import (
	"soloos/solodb/offheap"
	"sync"
)

type SrpcClientDriver struct {
	offheapDriver      *offheap.OffheapDriver
	netDriver          *NetDriver
	netConnReadSigPool offheap.RawObjectPool
	clientRWMutex      sync.RWMutex
	clients            map[PeerID]*SrpcClient
}

func (p *SrpcClientDriver) Init(offheapDriver *offheap.OffheapDriver, netDriver *NetDriver) error {
	var err error

	p.offheapDriver = offheapDriver
	p.netDriver = netDriver
	err = p.netConnReadSigPool.Init(int(offheap.MutexStrutSize), -1, nil, nil)
	if err != nil {
		return err
	}

	p.clients = make(map[PeerID]*SrpcClient)

	return nil
}

func (p *SrpcClientDriver) getClient(peerID PeerID) (*SrpcClient, error) {
	var (
		ret  *SrpcClient
		peer Peer
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

func (p *SrpcClientDriver) registerClient(peer Peer) (*SrpcClient, error) {
	var (
		client *SrpcClient
		err    error
	)

	p.clientRWMutex.Lock()
	client = p.clients[peer.ID]
	if client != nil {
		goto GET_CLIENT_DONE
	}

	client = &SrpcClient{}
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

func (p *SrpcClientDriver) CloseClient(peerID PeerID) error {
	var (
		client = p.clients[peerID]
		err    error
	)

	if client == nil {
		return nil
	}

	err = p.sendCloseCmd(client)
	if err != nil {
		return err
	}

	err = client.Close(ErrClosedByUser)
	if err != nil {
		return err
	}

	return nil
}
