package srpc

import (
	"soloos/common/snet/types"
	"soloos/sdbone/offheap"
)

func (p *Client) prepareWaitResponse(reqID uint64, resp *types.Response) error {
	resp.NetConnReadSig = offheap.MutexUintptr(p.clientDriver.netConnReadSigPool.AllocRawObject())
	resp.NetConnReadSig.Ptr().Lock()

	p.reqSigMapMutex.Lock()
	p.reqSigMap[reqID] = resp.NetConnReadSig
	p.reqSigMapMutex.Unlock()
	return nil
}

func (p *Client) activiateRequestSig(netQuery *types.NetQuery) error {
	var netConnReadSig offheap.MutexUintptr

	p.reqSigMapMutex.Lock()
	netConnReadSig = p.reqSigMap[netQuery.ReqID]
	delete(p.reqSigMap, netQuery.ReqID)
	p.reqSigMapMutex.Unlock()

	if netConnReadSig == 0 {
		return types.ErrObjectNotExists
	}

	p.doingNetQueryChan <- *netQuery

	// activiate req
	netConnReadSig.Ptr().Unlock()

	return nil
}

func (p *Client) doWaitResponse(req *types.Request, resp *types.Response) error {
	// wait cronReadResponse fetch data
	resp.NetConnReadSig.Ptr().Lock()
	resp.NetConnReadSig.Ptr().Unlock()

	p.clientDriver.netConnReadSigPool.ReleaseRawObject(uintptr(resp.NetConnReadSig))
	resp.NetConnReadSig = 0

	resp.NetQuery = <-p.doingNetQueryChan

	return nil
}
