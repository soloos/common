package srpc

import (
	"soloos/common/snettypes"
)

func (p *ClientDriver) sendCloseCmd(client *Client) error {
	var (
		req snettypes.Request
		err error
	)

	req.Init(client.AllocRequestID(), &client.doingNetQueryConn, "/Close")
	err = client.Write(&req)
	if err != nil {
		return err
	}

	return nil
}

func (p *ClientDriver) Call(uPeer snettypes.PeerUintptr,
	serviceID string,
	req *snettypes.Request,
	resp *snettypes.Response) error {
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

func (p *ClientDriver) AsyncCall(uPeer snettypes.PeerUintptr,
	serviceID string,
	req *snettypes.Request,
	resp *snettypes.Response) error {
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

func (p *ClientDriver) WaitResponse(uPeer snettypes.PeerUintptr,
	req *snettypes.Request,
	resp *snettypes.Response) error {
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
