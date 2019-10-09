package snet

func (p *SrpcClientDriver) ReadResponse(peerID PeerID,
	snetReq *SNetReq, snetResp *SNetResp, respBody []byte,
) error {
	var (
		client *SrpcClient
		err    error
	)

	client, err = p.getClient(peerID)
	if err != nil {
		return err
	}

	err = client.ReadResponse(snetReq, snetResp, respBody)
	if err != nil {
		return err
	}

	return err
}

func (p *SrpcClientDriver) SimpleReadResponse(peerID PeerID,
	snetReq *SNetReq, snetResp *SNetResp, resp IResponse,
	// snetReq *SNetReq, snetResp *SNetResp, resp *Response,
) error {
	var (
		client *SrpcClient
		err    error
	)

	client, err = p.getClient(peerID)
	if err != nil {
		return err
	}

	err = client.SimpleReadResponse(snetReq, snetResp, resp)
	if err != nil {
		return err
	}

	return err
}
