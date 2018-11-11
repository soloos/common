package snet

import (
	"soloos/snet/srpc"
	"soloos/snet/types"
	"soloos/util/offheap"
)

type ClientDriver struct {
	offheapDriver    *offheap.OffheapDriver
	rpcClientDrivers map[string]types.RpcClientDriver // map protocol to RpcClientDriver
}

func (p *ClientDriver) Init(offheapDriver *offheap.OffheapDriver) error {
	var err error
	p.offheapDriver = offheapDriver
	if err != nil {
		return err
	}

	p.rpcClientDrivers = make(map[string]types.RpcClientDriver)

	var srpcClientDriver = new(srpc.ClientDriver)
	err = srpcClientDriver.Init(p.offheapDriver)
	if err != nil {
		return err
	}
	p.rpcClientDrivers["srpc"] = srpcClientDriver

	return nil
}

func (p *ClientDriver) RegisterPeer(uPeer types.PeerUintptr) error {
	var (
		rpcClientDriver types.RpcClientDriver
		err             error
	)
	rpcClientDriver = p.rpcClientDrivers[uPeer.Ptr().ServiceProtocolStr()]

	err = rpcClientDriver.RegisterClient(uPeer)
	if err != nil {
		return err
	}

	return nil
}

func (p *ClientDriver) Call(uPeer types.PeerUintptr,
	serviceID string,
	request *types.Request,
	response *types.Response) error {
	var (
		rpcClientDriver types.RpcClientDriver
		err             error
	)
	rpcClientDriver = p.rpcClientDrivers[uPeer.Ptr().ServiceProtocolStr()]
	err = rpcClientDriver.Call(uPeer, serviceID, request, response)
	if err != nil {
		return err
	}

	return nil
}
