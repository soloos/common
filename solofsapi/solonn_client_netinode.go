package solofsapi

import (
	"soloos/common/solofsapitypes"
	"soloos/common/solofsprotocol"
)

func (p *SolonnClient) doGetNetINodeMetaData(isMustGet bool,
	uNetINode solofsapitypes.NetINodeUintptr,
	size uint64, netBlockCap int, memBlockCap int,
) error {
	var (
		req  solofsprotocol.NetINodeInfoReq
		resp solofsprotocol.NetINodeInfoResp
		err  error
	)

	req.NetINodeID = uNetINode.Ptr().ID
	req.Size = size
	req.NetBlockCap = int32(netBlockCap)
	req.MemBlockCap = int32(memBlockCap)

	if isMustGet {
		err = p.SNetClientDriver.SimpleCall(p.solonnPeerID,
			"/NetINode/MustGet", &resp, &req,
		)
	} else {
		err = p.SNetClientDriver.SimpleCall(p.solonnPeerID,
			"/NetINode/Get", &resp, &req,
		)
	}
	if err != nil {
		return err
	}

	uNetINode.Ptr().Size = resp.Size
	uNetINode.Ptr().NetBlockCap = int(resp.NetBlockCap)
	uNetINode.Ptr().MemBlockCap = int(resp.MemBlockCap)

	return nil
}

func (p *SolonnClient) GetNetINodeMetaData(uNetINode solofsapitypes.NetINodeUintptr) error {
	return p.doGetNetINodeMetaData(false, uNetINode, 0, 0, 0)
}

func (p *SolonnClient) MustGetNetINodeMetaData(uNetINode solofsapitypes.NetINodeUintptr,
	size uint64, netBlockCap int, memBlockCap int,
) error {
	return p.doGetNetINodeMetaData(true, uNetINode, size, netBlockCap, memBlockCap)
}

func (p *SolonnClient) NetINodeCommitSizeInDB(uNetINode solofsapitypes.NetINodeUintptr,
	size uint64) error {
	var req = solofsprotocol.NetINodeCommitSizeInDBReq{
		NetINodeID: uNetINode.Ptr().ID,
		Size:       size,
	}

	return p.SNetClientDriver.SimpleCall(p.solonnPeerID,
		"/NetINode/CommitSizeInDB", nil, req)
}
