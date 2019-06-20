package sdfsapi

import (
	"soloos/common/sdfsapitypes"
	"soloos/common/sdfsprotocol"
	"soloos/common/snettypes"

	flatbuffers "github.com/google/flatbuffers/go"
)

func (p *DataNodeClient) NetINodeSync(dataNodePeerID snettypes.PeerID,
	uNetINode sdfsapitypes.NetINodeUintptr) error {
	var (
		req             snettypes.Request
		resp            snettypes.Response
		protocolBuilder flatbuffers.Builder
		netINodeIDOff   flatbuffers.UOffsetT
		err             error
	)

	netINodeIDOff = protocolBuilder.CreateByteString(uNetINode.Ptr().ID[:])
	sdfsprotocol.NetINodeSyncRequestStart(&protocolBuilder)
	sdfsprotocol.NetINodeSyncRequestAddNetINodeID(&protocolBuilder, netINodeIDOff)
	protocolBuilder.Finish(sdfsprotocol.NetINodeSyncRequestEnd(&protocolBuilder))
	req.Param = protocolBuilder.Bytes[protocolBuilder.Head():]

	err = p.SNetClientDriver.Call(dataNodePeerID,
		"/NetINode/Sync", &req, &resp)
	if err != nil {
		return err
	}

	var body = make([]byte, resp.BodySize)[:resp.BodySize]
	err = p.SNetClientDriver.ReadResponse(dataNodePeerID, &req, &resp, body)
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
