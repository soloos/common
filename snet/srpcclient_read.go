package snet

func (p *SrpcClient) WaitResponse(snetReq *SNetReq, snetResp *SNetResp) error {
	return p.doWaitResponse(snetReq, snetResp)
}

func (p *SrpcClient) ReadResponse(snetReq *SNetReq, snetResp *SNetResp, respBody []byte) error {
	return snetResp.ReadAll(respBody)
}

func (p *SrpcClient) SimpleReadResponse(snetReq *SNetReq, snetResp *SNetResp, resp IResponse) error {
	// func (p *SrpcClient) SimpleReadResponse(snetReq *SNetReq, snetResp *SNetResp, resp *Response) error {
	return p.doingNetQueryConn.SimpleUnmarshalResponse(snetResp, resp)
}

func (p *SrpcClient) cronReadResponse() error {
	var (
		netQuery   NetQuery
		respHeader SNetRespHeader
		err        error
	)
	netQuery.Init(&p.doingNetQueryConn)

	for {
		err = netQuery.ReadSNetRespHeader(p.MaxMessageLength, &respHeader)
		if err != nil {
			break
			// goto FETCH_DATA_DONE
		}

		err = p.activiateRequestSig(&netQuery)
		if err != nil {
			break
			// goto FETCH_DATA_DONE
		}

		p.doingNetQueryConn.WaitReadDone()
	}

	// FETCH_DATA_DONE:
	if err != nil {
		err = p.doingNetQueryConn.Close(err)
	}
	return err
}
