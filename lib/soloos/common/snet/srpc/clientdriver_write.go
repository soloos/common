package srpc

import (
	"soloos/common/snet/types"
)

func (p *ClientDriver) sendCloseCmd(client *Client) error {
	var (
		req types.Request
		err error
	)

	req.Init(client.AllocRequestID(), &client.doingNetQueryConn, "/Close")
	err = client.Write(&req)
	if err != nil {
		return err
	}

	return nil
}

func (p *ClientDriver) Call(uPeer types.PeerUintptr,
	serviceID string,
	req *types.Request,
	resp *types.Response) error {
	var (
		err error
	)

	err = p.AsyncCall(uPeer, serviceID, req, resp)
	if err != nil {
		return err
	}

	err = p.WaitResponse(uPeer, req, resp)
	if err != nil {
		return err
	}

	return nil
}

func (p *ClientDriver) AsyncCall(uPeer types.PeerUintptr,
	serviceID string,
	req *types.Request,
	resp *types.Response) error {
	var (
		client *Client
		err    error
	)

	client, err = p.getClient(uPeer)
	if err != nil {
		return err
	}

	req.Init(client.AllocRequestID(), &client.doingNetQueryConn, serviceID)
	err = client.prepareWaitResponse(req.ReqID, resp)
	if err != nil {
		return err
	}

	err = client.Write(req)
	if err != nil {
		return err
	}

	return nil
}

func (p *ClientDriver) WaitResponse(uPeer types.PeerUintptr,
	req *types.Request,
	resp *types.Response) error {
	var (
		client *Client
		err    error
	)

	client, err = p.getClient(uPeer)
	if err != nil {
		return err
	}

	err = client.WaitResponse(req, resp)
	if err != nil {
		return err
	}

	return nil
}
