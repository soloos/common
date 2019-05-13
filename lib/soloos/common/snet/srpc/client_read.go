package srpc

import (
	"soloos/common/snet/types"
)

func (p *Client) WaitResponse(req *types.Request, resp *types.Response) error {
	return p.doWaitResponse(req, resp)
}

func (p *Client) ReadResponse(resp *types.Response, respBody []byte) error {
	return resp.NetQuery.ReadAll(respBody)
}

func (p *Client) cronReadResponse() error {
	var (
		netQuery   types.NetQuery
		respHeader types.ResponseHeader
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
