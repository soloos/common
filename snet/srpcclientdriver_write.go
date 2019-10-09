package snet

import (
	"soloos/common/log"
	"soloos/common/util"
)

func (p *SrpcClientDriver) sendCloseCmd(client *SrpcClient) error {
	var (
		snetReq SNetReq
		err     error
	)

	snetReq.Init(client.AllocRequestID(), &client.doingNetQueryConn, "/Close")
	err = client.Write(&snetReq)
	if err != nil {
		return err
	}

	return nil
}

func (p *SrpcClientDriver) Call(peerID PeerID,
	url string, snetReq *SNetReq, snetResp *SNetResp, reqArgs ...interface{},
) error {
	var (
		err error
	)

	err = p.AsyncCall(peerID, url, snetReq, snetResp, reqArgs...)
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
	url string, snetReq *SNetReq, snetResp *SNetResp, reqArgs ...interface{},
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

	if len(reqArgs) == 0 {
		snetReq.Param = snetReq.Param[:0]
		err = client.Write(snetReq)
		if err != nil {
			return err
		}
		return nil
	}

	var buf = snetCodecBytesPool.Get().(util.Buffer)
	defer snetCodecBytesPoolPut(buf)
	if len(reqArgs) == 1 {
		err = client.doingNetQueryConn.Marshal(&buf, reqArgs[0])
	} else {
		err = client.doingNetQueryConn.Marshal(&buf, reqArgs)
	}
	if err != nil {
		return err
	}

	snetReq.Param = buf.Bytes()
	err = client.Write(snetReq)
	if err != nil {
		return err
	}

	return nil
}

func (p *SrpcClientDriver) WaitResponse(peerID PeerID,
	snetReq *SNetReq, snetResp *SNetResp) error {
	var (
		client *SrpcClient
		err    error
	)

	client, err = p.getClient(peerID)
	if err != nil {
		return err
	}

	err = client.WaitResponse(snetReq, snetResp)
	if err != nil {
		return err
	}

	return nil
}

func (p *SrpcClientDriver) SimpleCall(peerID PeerID,
	url string, resp IResponse, reqArgs ...interface{},
	// url string, resp *Response, reqArgs ...interface{},
) error {
	var (
		snetReq  SNetReq
		snetResp SNetResp
		err      error
	)

	err = p.Call(peerID, url, &snetReq, &snetResp, reqArgs...)
	if err != nil {
		return err
	}

	err = p.SimpleReadResponse(peerID,
		&snetReq, &snetResp, resp,
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
