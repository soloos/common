package srpc

import (
	"soloos/snet/types"
	"soloos/util/offheap"
	"sync"
)

type ClientDriver struct {
	offheapDriver      *offheap.OffheapDriver
	netConnReadSigPool offheap.RawObjectPool
	clientRWMutex      sync.RWMutex
	clients            map[types.PeerUintptr]*Client
}

var _ = types.RpcClientDriver(&ClientDriver{})

func (p *ClientDriver) Init(offheapDriver *offheap.OffheapDriver) error {
	var err error

	p.offheapDriver = offheapDriver
	err = p.offheapDriver.InitRawObjectPool(&p.netConnReadSigPool, int(offheap.MutexStrutSize), -1, nil, nil)
	if err != nil {
		return err
	}

	p.clients = make(map[types.PeerUintptr]*Client)

	return nil
}

func (p *ClientDriver) setClient(uPeer types.PeerUintptr, client *Client) {
	p.clientRWMutex.Lock()
	p.clients[uPeer] = client
	p.clientRWMutex.Unlock()
}

func (p *ClientDriver) getClient(uPeer types.PeerUintptr) (ret *Client, err error) {
	p.clientRWMutex.RLock()
	ret = p.clients[uPeer]
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

	p.clients[uPeer] = ret

GET_CLIENT_DONE:
	p.clientRWMutex.Unlock()
	return
}

func (p *ClientDriver) sendCloseCmd(client *Client) error {
	var (
		req types.Request
		err error
	)

	err = client.Write(client.AllocRequestID(), "/Close", &req)
	if err != nil {
		return err
	}

	return nil
}

func (p *ClientDriver) RegisterClient(uPeer types.PeerUintptr, client interface{}) error {
	p.setClient(uPeer, client.(*Client))
	return nil
}

func (p *ClientDriver) CloseClient(uPeer types.PeerUintptr) error {
	var (
		client = p.clients[uPeer]
		err    error
	)

	err = p.sendCloseCmd(client)
	if err != nil {
		return err
	}

	err = client.Close()
	if err != nil {
		return err
	}

	return nil
}
