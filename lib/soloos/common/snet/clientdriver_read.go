package snet

import "soloos/common/snet/types"

func (p *ClientDriver) ReadResponse(uPeer types.PeerUintptr,
	req *types.Request,
	resp *types.Response,
	respBody []byte) error {
	var (
		rpcClientDriver types.RpcClientDriver
		err             error
	)
	rpcClientDriver = p.rpcClientDrivers[uPeer.Ptr().ServiceProtocol]
	err = rpcClientDriver.ReadResponse(uPeer, req, resp, respBody)
	if err != nil {
		return err
	}

	return nil
}
