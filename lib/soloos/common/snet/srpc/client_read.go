package srpc

import (
	"soloos/common/snettypes"
)

func (p *Client) WaitResponse(req *snettypes.Request, resp *snettypes.Response) error {
	return p.doWaitResponse(req, resp)
}

func (p *Client) ReadResponse(resp *snettypes.Response, respBody []byte) error {
	return resp.NetQuery.ReadAll(respBody)
}

func (p *Client) cronReadResponse() error {
	var (
		netQuery   snettypes.NetQuery
		respHeader snettypes.ResponseHeader
		err        error
	)

	netQuery.Init(&p.doingNetQueryConn)

	for {
		// fetch data
		err = netQuery.ReadResponseHeader(p.MaxMessageLength, &respHeader)
		if err != nil {
			goto FETCH_DATA_DONE
		}

		err = p.activiateRequestSig(&netQuery)
		if err != nil {
			goto FETCH_DATA_DONE
		}

		p.doingNetQueryConn.WaitReadDone()
	}

FETCH_DATA_DONE:
	if err != nil {
		err = p.doingNetQueryConn.Close(err)
	}
	return err
}
