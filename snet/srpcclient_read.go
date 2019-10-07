package snet

func (p *SrpcClient) WaitResponse(req *SNetReq, resp *SNetResp) error {
	return p.doWaitResponse(req, resp)
}

func (p *SrpcClient) ReadResponse(resp *SNetResp, respBody []byte) error {
	return resp.NetQuery.ReadAll(respBody)
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
