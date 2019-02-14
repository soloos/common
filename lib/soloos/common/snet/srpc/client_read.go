package srpc

import (
	"soloos/common/snet/types"
	"soloos/sdbone/offheap"
)

func (p *Client) PrepareWaitResponse(reqID uint64, resp *types.Response) error {
	resp.NetConnReadSig = offheap.MutexUintptr(p.clientDriver.netConnReadSigPool.AllocRawObject())
	resp.NetConnReadSig.Ptr().Lock()

	p.reqSigMapMutex.Lock()
	p.reqSigMap[reqID] = resp.NetConnReadSig
	p.reqSigMapMutex.Unlock()
	return nil
}

func (p *Client) WaitResponse(req *types.Request, resp *types.Response) error {
	// wait cronReadResponse fetch data
	resp.NetConnReadSig.Ptr().Lock()
	resp.NetConnReadSig.Ptr().Unlock()

	p.clientDriver.netConnReadSigPool.ReleaseRawObject(uintptr(resp.NetConnReadSig))
	resp.NetConnReadSig = 0

	// cronReadResponse fetched data
	resp.BodySize = p.lastResponseHeader.BodySize()
	resp.ParamSize = p.lastResponseHeader.ParamSize()

	return nil
}

func (p *Client) ReadResponse(respBody []byte) error {
	return p.Conn.ReadAll(respBody)
}

func (p *Client) activiateRequestSig(reqID uint64) error {
	var netConnReadSig offheap.MutexUintptr

	p.reqSigMapMutex.Lock()
	netConnReadSig = p.reqSigMap[reqID]
	delete(p.reqSigMap, reqID)
	p.reqSigMapMutex.Unlock()

	// activiate req
	netConnReadSig.Ptr().Unlock()

	return nil
}

func (p *Client) cronReadResponse() error {
	var (
		reqID uint64
		err   error
	)

	for {
		// fetch data
		err = p.Conn.ReadResponseHeader(p.MaxMessageLength, &p.lastResponseHeader)
		if err != nil {
			goto FETCH_DATA_DONE
		}

		reqID = p.lastResponseHeader.ID()

		err = p.Conn.AfterReadHeaderSuccess()
		if err != nil {
			goto FETCH_DATA_DONE
		}

		err = p.activiateRequestSig(reqID)
		if err != nil {
			goto FETCH_DATA_DONE
		}

		// wait read done
		p.Conn.ContinueReadSig.Lock()
		p.Conn.ContinueReadSig.Unlock()
	}

FETCH_DATA_DONE:
	if err != nil {
		err = p.Conn.Close(err)
	}
	return err
}
