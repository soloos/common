package solofsapi

import (
	"soloos/common/solofsapitypes"
	"soloos/common/solofsprotocol"
	"soloos/common/snettypes"

	flatbuffers "github.com/google/flatbuffers/go"
)

func (p *SolonnClient) PrepareNetBlockMetaData(netBlockInfo *solofsprotocol.NetINodeNetBlockInfoResponse,
	uNetINode solofsapitypes.NetINodeUintptr,
	netBlockIndex int32,
	uNetBlock solofsapitypes.NetBlockUintptr,
) error {
	var (
		req             snettypes.Request
		resp            snettypes.Response
		protocolBuilder flatbuffers.Builder
		netINodeIDOff   flatbuffers.UOffsetT
		err             error
	)

	netINodeIDOff = protocolBuilder.CreateString(uNetINode.Ptr().IDStr())
	solofsprotocol.NetINodeNetBlockInfoRequestStart(&protocolBuilder)
	solofsprotocol.NetINodeNetBlockInfoRequestAddNetINodeID(&protocolBuilder, netINodeIDOff)
	solofsprotocol.NetINodeNetBlockInfoRequestAddNetBlockIndex(&protocolBuilder, int32(netBlockIndex))
	solofsprotocol.NetINodeNetBlockInfoRequestAddCap(&protocolBuilder, int32(uNetINode.Ptr().NetBlockCap))
	protocolBuilder.Finish(solofsprotocol.NetINodeNetBlockInfoRequestEnd(&protocolBuilder))
	req.Param = protocolBuilder.Bytes[protocolBuilder.Head():]

	// TODO choose solonn
	err = p.SNetClientDriver.Call(p.solonnPeerID,
		"/NetBlock/PrepareMetaData", &req, &resp)
	if err != nil {
		return err
	}

	var body = make([]byte, resp.BodySize)[:resp.BodySize]
	err = p.SNetClientDriver.ReadResponse(p.solonnPeerID, &req, &resp, body)
	if err != nil {
		return err
	}

	var (
		pNetBlock      = uNetBlock.Ptr()
		commonResponse solofsprotocol.CommonResponse
	)
	netBlockInfo.Init(body, flatbuffers.GetUOffsetT(body))
	netBlockInfo.CommonResponse(&commonResponse)
	err = CommonResponseToError(&commonResponse)
	if err != nil {
		return err
	}

	pNetBlock.NetINodeID = uNetINode.Ptr().ID
	pNetBlock.IndexInNetINode = netBlockIndex
	pNetBlock.Len = int(netBlockInfo.Len())
	pNetBlock.Cap = int(netBlockInfo.Cap())

	return nil
}
