package snet

import (
	"soloos/solodb/offheap"
)

func (p *SrpcClient) prepareWaitResponse(reqID uint64, resp *SNetResp) error {
	resp.NetConnReadSig = offheap.MutexUintptr(p.clientDriver.netConnReadSigPool.AllocRawObject())
	resp.NetConnReadSig.Ptr().Lock()

	p.reqSigMapMutex.Lock()
	p.reqSigMap[reqID] = resp.NetConnReadSig
	p.reqSigMapMutex.Unlock()
	return nil
}

func (p *SrpcClient) activiateRequestSig(netQuery *NetQuery) error {
	var netConnReadSig offheap.MutexUintptr

	p.reqSigMapMutex.Lock()
	netConnReadSig = p.reqSigMap[netQuery.ReqID]
	delete(p.reqSigMap, netQuery.ReqID)
	p.reqSigMapMutex.Unlock()

	if netConnReadSig == 0 {
		return ErrObjectNotExists
	}

	p.doingNetQueryChan <- *netQuery

	// activiate req
	netConnReadSig.Ptr().Unlock()

	return nil
}

func (p *SrpcClient) doWaitResponse(snetReq *SNetReq, snetResp *SNetResp) error {
	// wait cronReadResponse fetch data
	snetResp.NetConnReadSig.Ptr().Lock()
	snetResp.NetConnReadSig.Ptr().Unlock()

	p.clientDriver.netConnReadSigPool.ReleaseRawObject(uintptr(snetResp.NetConnReadSig))
	snetResp.NetConnReadSig = 0

	snetResp.NetQuery = <-p.doingNetQueryChan

	return nil
}
