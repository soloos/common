package main

import (
	"soloos/common/snet"
	"soloos/common/snet/types"
	"soloos/common/util"
	"soloos/solodb/offheap"
)

type Tool struct {
	offheapDriver    *offheap.OffheapDriver
	SNetDriver       snet.NetDriver
	SNetClientDriver snet.ClientDriver
}

func (p *Tool) Init() error {
	var err error

	p.offheapDriver = &offheap.DefaultOffheapDriver

	err = p.SNetDriver.Init(p.offheapDriver, "Tool")
	if err != nil {
		return err
	}

	err = p.SNetClientDriver.Init(p.offheapDriver)
	if err != nil {
		return err
	}

	return nil
}

func (p *Tool) Test() error {
	var (
		uPeer snettypes.PeerUintptr
		req   snettypes.Request
		resp  snettypes.Response
		err   error
	)
	uPeer = p.SNetDriver.AllocPeer("127.0.0.1:1339", snettypes.ProtocolSRPC)
	err = p.SNetClientDriver.Call(uPeer, "/Test", &req, &resp)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	var tool Tool
	util.AssertErrIsNil(tool.Init())
	util.AssertErrIsNil(tool.Test())

}
