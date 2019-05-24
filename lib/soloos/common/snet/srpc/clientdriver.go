package srpc

import (
	"soloos/common/snettypes"
	"soloos/sdbone/offheap"
	"sync"
)

type ClientDriver struct {
	offheapDriver      *offheap.OffheapDriver
	netConnReadSigPool offheap.RawObjectPool
	clientRWMutex      sync.RWMutex
	clients            map[snettypes.PeerID]*Client
}

var _ = snettypes.RpcClientDriver(&ClientDriver{})

func (p *ClientDriver) Init(offheapDriver *offheap.OffheapDriver) error {
	var err error

	p.offheapDriver = offheapDriver
	err = p.netConnReadSigPool.Init(int(offheap.MutexStrutSize), -1, nil, nil)
	if err != nil {
		return err
	}

	p.clients = make(map[snettypes.PeerID]*Client)

	return nil
}

func (p *ClientDriver) setClient(uPeer snettypes.PeerUintptr, client *Client) {
	p.clientRWMutex.Lock()
	p.clients[uPeer.Ptr().ID] = client
	p.clientRWMutex.Unlock()
}

func (p *ClientDriver) getClient(uPeer snettypes.PeerUintptr) (ret *Client, err error) {
	p.clientRWMutex.RLock()
	ret = p.clients[uPeer.Ptr().ID]
	p.clientRWMutex.RUnlock()

	if ret != nil {
		return
	}

	p.clientRWMutex.Lock()
	if ret != nil {
		goto GET_CLIENT_DONE
	}

	ret = &Client{}
	err = ret.Init(p, uPeer.Ptr().AddressStr())
	if err != nil {
		ret = nil
		goto GET_CLIENT_DONE
	}

	err = ret.Start()
	if err != nil {
		ret = nil
		goto GET_CLIENT_DONE
	}

	p.clients[uPeer.Ptr().ID] = ret

GET_CLIENT_DONE:
	p.clientRWMutex.Unlock()
	return
}

func (p *ClientDriver) RegisterClient(uPeer snettypes.PeerUintptr, client interface{}) error {
	p.setClient(uPeer, client.(*Client))
	return nil
}

func (p *ClientDriver) CloseClient(uPeer snettypes.PeerUintptr) error {
	var (
		client = p.clients[uPeer.Ptr().ID]
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
