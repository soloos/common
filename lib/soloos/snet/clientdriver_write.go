package snet

import "soloos/snet/types"

func (p *ClientDriver) Call(uPeer types.PeerUintptr,
	serviceID string,
	req *types.Request,
	resp *types.Response) error {
	var (
		rpcClientDriver types.RpcClientDriver
		err             error
	)
	rpcClientDriver = p.rpcClientDrivers[uPeer.Ptr().ServiceProtocol]
	err = rpcClientDriver.Call(uPeer, serviceID, req, resp)
	if err != nil {
		return err
	}

	return nil
}
