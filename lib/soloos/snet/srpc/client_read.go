package srpc

import (
	"soloos/snet/types"
	"soloos/util/offheap"
)

func (p *Client) PrepareWaitResponse(requestID uint64, response *types.Response) error {
	response.NetConnReadSig = offheap.MutexUintptr(p.clientDriver.netConnReadSigPool.AllocRawObject())
	response.NetConnReadSig.Ptr().Lock()

	p.requestSigMapMutex.Lock()
	p.requestSigMap[requestID] = response.NetConnReadSig
	p.requestSigMapMutex.Unlock()
	return nil
}

func (p *Client) ActiviateRequestSig(requestID uint64) error {
	var netConnReadSig offheap.MutexUintptr

	p.requestSigMapMutex.Lock()
	netConnReadSig = p.requestSigMap[requestID]
	delete(p.requestSigMap, requestID)
	p.requestSigMapMutex.Unlock()

	// activiate request
	netConnReadSig.Ptr().Unlock()

	return nil
}

func (p *Client) WaitResponse(request *types.Request, response *types.Response) error {
	// wait cronReadResponse fetch data
	response.NetConnReadSig.Ptr().Lock()
	response.NetConnReadSig.Ptr().Unlock()

	p.clientDriver.netConnReadSigPool.ReleaseRawObject(uintptr(response.NetConnReadSig))

	// cronReadResponse fetched data
	response.BodySize = p.lastResponseHeader.ContentLen()

	return nil
}

func (p *Client) ReadResponse(resp *[]byte) error {
	return p.Conn.ReadAll(*resp)
}

func (p *Client) doCronReadResponse() {
	var (
		requestID uint64
		err       error
	)

	for {
		// fetch data
		err = p.Conn.ReadResponseHeader(p.MaxMessageLength, &p.lastResponseHeader)
		if err != nil {
			goto FETCH_DATA_DONE
		}

		requestID = p.lastResponseHeader.ID()

		err = p.Conn.AfterReadHeaderSuccess()
		if err != nil {
			goto FETCH_DATA_DONE
		}

		err = p.ActiviateRequestSig(requestID)
		if err != nil {
			panic(err)
			goto FETCH_DATA_DONE
		}

		// wait read done
		p.Conn.ContinueReadSig.Lock()
		p.Conn.ContinueReadSig.Unlock()
	}

FETCH_DATA_DONE:
	if err != nil {
		// log.Warn("cronReadResponse", err)
	}
}

func (p *Client) cronReadResponse() error {
	go func() {
		p.doCronReadResponse()
	}()
	return nil
}
