package snet

import (
	"soloos/common/log"
	"soloos/common/snettypes"
)

func (p *SRPCClientDriver) sendCloseCmd(client *SRPCClient) error {
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

func (p *SRPCClientDriver) Call(peerID snettypes.PeerID,
	serviceID string, req *snettypes.Request, resp *snettypes.Response) error {
	var (
		err error
	)

	err = p.AsyncCall(peerID, serviceID, req, resp)
	if err != nil {
		log.Info("AsyncCall error ", err)
		return err
	}

	err = p.WaitResponse(peerID, req, resp)
	if err != nil {
		return err
	}

	return nil
}

func (p *SRPCClientDriver) AsyncCall(peerID snettypes.PeerID,
	serviceID string, req *snettypes.Request, resp *snettypes.Response) error {
	var (
		client *SRPCClient
		err    error
	)

	client, err = p.getClient(peerID)
	if err != nil {
		log.Info("AsyncCall getClient error ", err)
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

func (p *SRPCClientDriver) WaitResponse(peerID snettypes.PeerID,
	req *snettypes.Request, resp *snettypes.Response) error {
	var (
		client *SRPCClient
		err    error
	)

	client, err = p.getClient(peerID)
	if err != nil {
		return err
	}

	err = client.WaitResponse(req, resp)
	if err != nil {
		return err
	}

	return nil
}
