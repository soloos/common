package snet

import (
	"soloos/common/iron"
	"soloos/common/log"
	"soloos/common/snettypes"
)

func (p *SrpcClientDriver) sendCloseCmd(client *SrpcClient) error {
	var (
		req snettypes.SNetReq
		err error
	)

	req.Init(client.AllocRequestID(), &client.doingNetQueryConn, "/Close")
	err = client.Write(&req)
	if err != nil {
		return err
	}

	return nil
}

func (p *SrpcClientDriver) Call(peerID snettypes.PeerID,
	url string,
	snetReq *snettypes.SNetReq, snetResp *snettypes.SNetResp,
	req interface{},
) error {
	var (
		err error
	)

	err = p.AsyncCall(peerID, url, snetReq, snetResp, req)
	if err != nil {
		log.Info("AsyncCall error ", err)
		return err
	}

	err = p.WaitResponse(peerID, snetReq, snetResp)
	if err != nil {
		return err
	}

	return nil
}

func (p *SrpcClientDriver) AsyncCall(peerID snettypes.PeerID,
	url string,
	snetReq *snettypes.SNetReq, snetResp *snettypes.SNetResp,
	req interface{},
) error {
	var (
		client *SrpcClient
		err    error
	)

	client, err = p.getClient(peerID)
	if err != nil {
		log.Info("AsyncCall getClient error ", err)
		return err
	}

	snetReq.Init(client.AllocRequestID(), &client.doingNetQueryConn, url)
	snetReq.Param = iron.MustSpecMarshalRequest(req)

	err = client.prepareWaitResponse(snetReq.ReqID, snetResp)
	if err != nil {
		return err
	}

	err = client.Write(snetReq)
	if err != nil {
		return err
	}

	return nil
}

func (p *SrpcClientDriver) WaitResponse(peerID snettypes.PeerID,
	req *snettypes.SNetReq, resp *snettypes.SNetResp) error {
	var (
		client *SrpcClient
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

func (p *SrpcClientDriver) SimpleCall(peerID snettypes.PeerID,
	url string, req interface{}, ret interface{},
) error {
	var (
		snetReq  snettypes.SNetReq
		snetResp snettypes.SNetResp
		err      error
	)

	err = p.Call(peerID, url, &snetReq, &snetResp, req)
	if err != nil {
		return err
	}

	var snetRespBody = make([]byte, int(snetResp.ParamSize))
	err = p.ReadResponse(peerID, &snetReq, &snetResp,
		snetRespBody, ret,
	)
	if err != nil {
		return err
	}

	err = snetResp.SkipReadRemaining()
	if err != nil {
		return err
	}

	return nil
}
