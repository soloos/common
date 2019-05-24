package sdfsapi

import (
	"soloos/common/sdfsapitypes"
	"soloos/common/snettypes"
	"soloos/common/sdfsprotocol"

	flatbuffers "github.com/google/flatbuffers/go"
)

func (p *NameNodeClient) PrepareNetBlockMetaData(netBlockInfo *sdfsprotocol.NetINodeNetBlockInfoResponse,
	uNetINode sdfsapitypes.NetINodeUintptr,
	netBlockIndex int32,
	uNetBlock sdfsapitypes.NetBlockUintptr,
) error {
	var (
		req             snettypes.Request
		resp            snettypes.Response
		protocolBuilder flatbuffers.Builder
		netINodeIDOff   flatbuffers.UOffsetT
		err             error
	)

	netINodeIDOff = protocolBuilder.CreateString(uNetINode.Ptr().IDStr())
	sdfsprotocol.NetINodeNetBlockInfoRequestStart(&protocolBuilder)
	sdfsprotocol.NetINodeNetBlockInfoRequestAddNetINodeID(&protocolBuilder, netINodeIDOff)
	sdfsprotocol.NetINodeNetBlockInfoRequestAddNetBlockIndex(&protocolBuilder, int32(netBlockIndex))
	sdfsprotocol.NetINodeNetBlockInfoRequestAddCap(&protocolBuilder, int32(uNetINode.Ptr().NetBlockCap))
	protocolBuilder.Finish(sdfsprotocol.NetINodeNetBlockInfoRequestEnd(&protocolBuilder))
	req.Param = protocolBuilder.Bytes[protocolBuilder.Head():]

	// TODO choose namenode
	err = p.SNetClientDriver.Call(p.nameNodePeer,
		"/NetBlock/PrepareMetaData", &req, &resp)
	if err != nil {
		return err
	}

	var body = make([]byte, resp.BodySize)[:resp.BodySize]
	err = p.SNetClientDriver.ReadResponse(p.nameNodePeer, &req, &resp, body)
	if err != nil {
		return err
	}

	var (
		pNetBlock      = uNetBlock.Ptr()
		commonResponse sdfsprotocol.CommonResponse
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
