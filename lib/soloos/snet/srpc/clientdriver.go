package srpc

import (
	"soloos/snet/types"
	"soloos/util/offheap"
)

type ClientDriver struct {
	offheapDriver      *offheap.OffheapDriver
	netConnReadSigPool offheap.RawObjectPool
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

func (p *ClientDriver) RegisterClient(uPeer types.PeerUintptr) error {
	var (
		client = &Client{}
		err    error
	)
	err = client.Init(p, uPeer.Ptr().AddressStr())
	if err != nil {
		return err
	}

	err = client.Start()
	if err != nil {
		return err
	}

	p.clients[uPeer] = client

	return nil
}

func (p *ClientDriver) sendCloseCmd(client *Client) error {
	var (
		request types.Request
		err     error
	)

	err = client.Write(client.AllocRequestID(), "/Close", &request)
	if err != nil {
		return err
	}

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
