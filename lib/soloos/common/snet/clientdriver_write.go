package snet

import "soloos/common/snettypes"

func (p *ClientDriver) Call(uPeer snettypes.PeerUintptr,
	serviceID string,
	req *snettypes.Request,
	resp *snettypes.Response) error {
	var (
		rpcClientDriver snettypes.RpcClientDriver
		err             error
	)
	rpcClientDriver = p.rpcClientDrivers[uPeer.Ptr().ServiceProtocol]
	err = rpcClientDriver.Call(uPeer, serviceID, req, resp)
	if err != nil {
		return err
	}

	return nil
}
