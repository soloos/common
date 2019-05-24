package sdfsapi

import (
	"soloos/common/sdfsapitypes"
	"soloos/common/snettypes"
	"soloos/common/sdfsprotocol"

	flatbuffers "github.com/google/flatbuffers/go"
)

func (p *NameNodeClient) doGetNetINodeMetaData(isMustGet bool,
	uNetINode sdfsapitypes.NetINodeUintptr,
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
	sdfsprotocol.NetINodeInfoRequestStart(&protocolBuilder)
	sdfsprotocol.NetINodeInfoRequestAddNetINodeID(&protocolBuilder, netINodeIDOff)
	sdfsprotocol.NetINodeInfoRequestAddSize(&protocolBuilder, size)
	sdfsprotocol.NetINodeInfoRequestAddNetBlockCap(&protocolBuilder, int32(netBlockCap))
	sdfsprotocol.NetINodeInfoRequestAddMemBlockCap(&protocolBuilder, int32(memBlockCap))
	protocolBuilder.Finish(sdfsprotocol.NetINodeNetBlockInfoRequestEnd(&protocolBuilder))
	req.Param = protocolBuilder.Bytes[protocolBuilder.Head():]

	if isMustGet {
		err = p.SNetClientDriver.Call(p.nameNodePeer,
			"/NetINode/MustGet", &req, &resp)
	} else {
		err = p.SNetClientDriver.Call(p.nameNodePeer,
			"/NetINode/Get", &req, &resp)
	}
	if err != nil {
		return err
	}

	var body = make([]byte, resp.BodySize)[:resp.BodySize]
	err = p.SNetClientDriver.ReadResponse(p.nameNodePeer, &req, &resp, body)
	if err != nil {
		return err
	}

	var (
		netINodeInfo   sdfsprotocol.NetINodeInfoResponse
		commonResponse sdfsprotocol.CommonResponse
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

func (p *NameNodeClient) GetNetINodeMetaData(uNetINode sdfsapitypes.NetINodeUintptr) error {
	return p.doGetNetINodeMetaData(false, uNetINode, 0, 0, 0)
}

func (p *NameNodeClient) MustGetNetINodeMetaData(uNetINode sdfsapitypes.NetINodeUintptr,
	size uint64, netBlockCap int, memBlockCap int,
) error {
	return p.doGetNetINodeMetaData(true, uNetINode, size, netBlockCap, memBlockCap)
}

func (p *NameNodeClient) NetINodeCommitSizeInDB(uNetINode sdfsapitypes.NetINodeUintptr,
	size uint64) error {
	var (
		req             snettypes.Request
		resp            snettypes.Response
		protocolBuilder flatbuffers.Builder
		netINodeIDOff   flatbuffers.UOffsetT
		err             error
	)

	netINodeIDOff = protocolBuilder.CreateByteString(uNetINode.Ptr().ID[:])
	sdfsprotocol.NetINodeCommitSizeInDBRequestStart(&protocolBuilder)
	sdfsprotocol.NetINodeCommitSizeInDBRequestAddNetINodeID(&protocolBuilder, netINodeIDOff)
	sdfsprotocol.NetINodeCommitSizeInDBRequestAddSize(&protocolBuilder, size)
	protocolBuilder.Finish(sdfsprotocol.NetINodeCommitSizeInDBRequestEnd(&protocolBuilder))
	req.Param = protocolBuilder.Bytes[protocolBuilder.Head():]

	err = p.SNetClientDriver.Call(p.nameNodePeer,
		"/NetINode/CommitSizeInDB", &req, &resp)
	if err != nil {
		return err
	}

	var body = make([]byte, resp.BodySize)[:resp.BodySize]
	err = p.SNetClientDriver.ReadResponse(p.nameNodePeer, &req, &resp, body)
	if err != nil {
		return err
	}

	var (
		commonResponse sdfsprotocol.CommonResponse
	)

	commonResponse.Init(body, flatbuffers.GetUOffsetT(body))
	err = CommonResponseToError(&commonResponse)
	if err != nil {
		return err
	}

	return nil
}
