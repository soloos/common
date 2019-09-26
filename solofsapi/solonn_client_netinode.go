package solofsapi

import (
	"soloos/common/solofsapitypes"
	"soloos/common/solofsprotocol"
	"soloos/common/snettypes"

	flatbuffers "github.com/google/flatbuffers/go"
)

func (p *SolonnClient) doGetNetINodeMetaData(isMustGet bool,
	uNetINode solofsapitypes.NetINodeUintptr,
	size uint64, netBlockCap int, memBlockCap int,
) error {
	var (
		req             snettypes.Request
		resp            snettypes.Response
		protocolBuilder flatbuffers.Builder
		netINodeIDOff   flatbuffers.UOffsetT
		err             error
	)

	netINodeIDOff = protocolBuilder.CreateByteString(uNetINode.Ptr().ID[:])
	solofsprotocol.NetINodeInfoRequestStart(&protocolBuilder)
	solofsprotocol.NetINodeInfoRequestAddNetINodeID(&protocolBuilder, netINodeIDOff)
	solofsprotocol.NetINodeInfoRequestAddSize(&protocolBuilder, size)
	solofsprotocol.NetINodeInfoRequestAddNetBlockCap(&protocolBuilder, int32(netBlockCap))
	solofsprotocol.NetINodeInfoRequestAddMemBlockCap(&protocolBuilder, int32(memBlockCap))
	protocolBuilder.Finish(solofsprotocol.NetINodeNetBlockInfoRequestEnd(&protocolBuilder))
	req.Param = protocolBuilder.Bytes[protocolBuilder.Head():]

	if isMustGet {
		err = p.SNetClientDriver.Call(p.solonnPeerID,
			"/NetINode/MustGet", &req, &resp)
	} else {
		err = p.SNetClientDriver.Call(p.solonnPeerID,
			"/NetINode/Get", &req, &resp)
	}
	if err != nil {
		return err
	}

	var body = make([]byte, resp.BodySize)[:resp.BodySize]
	err = p.SNetClientDriver.ReadResponse(p.solonnPeerID, &req, &resp, body)
	if err != nil {
		return err
	}

	var (
		netINodeInfo   solofsprotocol.NetINodeInfoResponse
		commonResponse solofsprotocol.CommonResponse
	)

	netINodeInfo.Init(body, flatbuffers.GetUOffsetT(body))
	netINodeInfo.CommonResponse(&commonResponse)
	err = CommonResponseToError(&commonResponse)
	if err != nil {
		return err
	}

	uNetINode.Ptr().Size = netINodeInfo.Size()
	uNetINode.Ptr().NetBlockCap = int(netINodeInfo.NetBlockCap())
	uNetINode.Ptr().MemBlockCap = int(netINodeInfo.MemBlockCap())

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
	var (
		req             snettypes.Request
		resp            snettypes.Response
		protocolBuilder flatbuffers.Builder
		netINodeIDOff   flatbuffers.UOffsetT
		err             error
	)

	netINodeIDOff = protocolBuilder.CreateByteString(uNetINode.Ptr().ID[:])
	solofsprotocol.NetINodeCommitSizeInDBRequestStart(&protocolBuilder)
	solofsprotocol.NetINodeCommitSizeInDBRequestAddNetINodeID(&protocolBuilder, netINodeIDOff)
	solofsprotocol.NetINodeCommitSizeInDBRequestAddSize(&protocolBuilder, size)
	protocolBuilder.Finish(solofsprotocol.NetINodeCommitSizeInDBRequestEnd(&protocolBuilder))
	req.Param = protocolBuilder.Bytes[protocolBuilder.Head():]

	err = p.SNetClientDriver.Call(p.solonnPeerID,
		"/NetINode/CommitSizeInDB", &req, &resp)
	if err != nil {
		return err
	}

	var body = make([]byte, resp.BodySize)[:resp.BodySize]
	err = p.SNetClientDriver.ReadResponse(p.solonnPeerID, &req, &resp, body)
	if err != nil {
		return err
	}

	var (
		commonResponse solofsprotocol.CommonResponse
	)

	commonResponse.Init(body, flatbuffers.GetUOffsetT(body))
	err = CommonResponseToError(&commonResponse)
	if err != nil {
		return err
	}

	return nil
}
