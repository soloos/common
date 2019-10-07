package snet

import (
	"soloos/common/iron"
	"soloos/common/log"
)

func (p *SrpcClientDriver) sendCloseCmd(client *SrpcClient) error {
	var (
		req SNetReq
		err error
	)

	req.Init(client.AllocRequestID(), &client.doingNetQueryConn, "/Close")
	err = client.Write(&req)
	if err != nil {
		return err
	}

	return nil
}

func (p *SrpcClientDriver) Call(peerID PeerID,
	url string,
	snetReq *SNetReq, snetResp *SNetResp,
) error {
	var (
		err error
	)

	err = p.AsyncCall(peerID, url, snetReq, snetResp)
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

func (p *SrpcClientDriver) AsyncCall(peerID PeerID,
	url string,
	snetReq *SNetReq, snetResp *SNetResp,
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

func (p *SrpcClientDriver) WaitResponse(peerID PeerID,
	req *SNetReq, resp *SNetResp) error {
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

func (p *SrpcClientDriver) SimpleCall(peerID PeerID,
	url string, ret interface{}, reqArgs ...interface{},
) error {
	var (
		snetReq  SNetReq
		snetResp SNetResp
		err      error
	)

	if len(reqArgs) == 1 {
		snetReq.Param = iron.MustSpecMarshalRequest(reqArgs[0])
	} else {
		snetReq.Param = iron.MustSpecMarshalRequest(reqArgs)
	}
	err = p.Call(peerID, url, &snetReq, &snetResp)
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
