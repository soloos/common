package solofsapi

import (
	"soloos/common/snettypes"
	"soloos/common/solofsapitypes"
	"soloos/common/solofsprotocol"
)

func (p *SolonnClient) PrepareNetBlockMetaData(netBlockInfo *solofsprotocol.NetINodeNetBlockInfoResp,
	uNetINode solofsapitypes.NetINodeUintptr,
	netBlockIndex int32,
	uNetBlock solofsapitypes.NetBlockUintptr,
) error {
	var (
		snetReq  snettypes.SNetReq
		snetResp snettypes.SNetResp
		req      solofsprotocol.NetINodeNetBlockInfoReq
		err      error
	)

	req.NetINodeID = uNetINode.Ptr().ID
	req.NetBlockIndex = int32(netBlockIndex)
	req.Cap = int32(uNetINode.Ptr().NetBlockCap)

	// TODO choose solonn
	err = p.SNetClientDriver.Call(p.solonnPeerID,
		"/NetBlock/PrepareMetaData", &snetReq, &snetResp, req)
	if err != nil {
		return err
	}

	var (
		pNetBlock = uNetBlock.Ptr()
	)
	var respParamBs = make([]byte, snetResp.ParamSize)
	err = p.SNetClientDriver.ReadResponse(p.solonnPeerID, &snetReq, &snetResp,
		respParamBs, netBlockInfo)
	if err != nil {
		return err
	}

	pNetBlock.NetINodeID = uNetINode.Ptr().ID
	pNetBlock.IndexInNetINode = netBlockIndex
	pNetBlock.Len = int(netBlockInfo.Len)
	pNetBlock.Cap = int(netBlockInfo.Cap)

	return nil
}
