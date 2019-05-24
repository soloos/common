package snet

import "soloos/common/snettypes"

func (p *ClientDriver) ReadResponse(uPeer snettypes.PeerUintptr,
	req *snettypes.Request,
	resp *snettypes.Response,
	respBody []byte) error {
	var (
		rpcClientDriver snettypes.RpcClientDriver
		err             error
	)
	rpcClientDriver = p.rpcClientDrivers[uPeer.Ptr().ServiceProtocol]
	err = rpcClientDriver.ReadResponse(uPeer, req, resp, respBody)
	if err != nil {
		return err
	}

	return nil
}
