package snet

import (
	"soloos/common/snet/srpc"
	"soloos/common/snettypes"
	"soloos/sdbone/offheap"
)

type ClientDriver struct {
	offheapDriver    *offheap.OffheapDriver
	rpcClientDrivers [255]snettypes.RpcClientDriver // map protocol to RpcClientDriver
}

func (p *ClientDriver) Init(offheapDriver *offheap.OffheapDriver) error {
	var err error
	p.offheapDriver = offheapDriver
	if err != nil {
		return err
	}

	// p.rpcClientDrivers = make(map[int]RpcClientDriver)

	var srpcClientDriver = new(srpc.ClientDriver)
	err = srpcClientDriver.Init(p.offheapDriver)
	if err != nil {
		return err
	}
	p.rpcClientDrivers[snettypes.ProtocolSRPC] = srpcClientDriver

	return nil
}

func (p *ClientDriver) GetRPCClientDriver(protocol int) snettypes.RpcClientDriver {
	return p.rpcClientDrivers[protocol]
}
