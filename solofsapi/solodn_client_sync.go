package solofsapi

import (
	"soloos/common/solofsapitypes"
	"soloos/common/solofsprotocol"
	"soloos/common/snettypes"

	flatbuffers "github.com/google/flatbuffers/go"
)

func (p *SolodnClient) NetINodeSync(solodnPeerID snettypes.PeerID,
	uNetINode solofsapitypes.NetINodeUintptr) error {
	var (
		req             snettypes.Request
		resp            snettypes.Response
		protocolBuilder flatbuffers.Builder
		netINodeIDOff   flatbuffers.UOffsetT
		err             error
	)

	netINodeIDOff = protocolBuilder.CreateByteString(uNetINode.Ptr().ID[:])
	solofsprotocol.NetINodeSyncRequestStart(&protocolBuilder)
	solofsprotocol.NetINodeSyncRequestAddNetINodeID(&protocolBuilder, netINodeIDOff)
	protocolBuilder.Finish(solofsprotocol.NetINodeSyncRequestEnd(&protocolBuilder))
	req.Param = protocolBuilder.Bytes[protocolBuilder.Head():]

	err = p.SNetClientDriver.Call(solodnPeerID,
		"/NetINode/Sync", &req, &resp)
	if err != nil {
		return err
	}

	var body = make([]byte, resp.BodySize)[:resp.BodySize]
	err = p.SNetClientDriver.ReadResponse(solodnPeerID, &req, &resp, body)
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
