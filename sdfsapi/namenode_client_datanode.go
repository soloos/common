package sdfsapi

import (
	"soloos/common/sdfsprotocol"
	"soloos/common/snettypes"

	flatbuffers "github.com/google/flatbuffers/go"
)

func (p *NameNodeClient) RegisterDataNode(peerID snettypes.PeerID,
	serveAddr string,
	protocolType int) error {
	var (
		req             snettypes.Request
		resp            snettypes.Response
		protocolBuilder flatbuffers.Builder
		peerIDOff       flatbuffers.UOffsetT
		addrOff         flatbuffers.UOffsetT
		err             error
	)

	peerIDOff = protocolBuilder.CreateByteString(peerID[:])
	addrOff = protocolBuilder.CreateString(serveAddr)
	sdfsprotocol.SNetPeerStart(&protocolBuilder)
	sdfsprotocol.SNetPeerAddPeerID(&protocolBuilder, peerIDOff)
	sdfsprotocol.SNetPeerAddAddress(&protocolBuilder, addrOff)
	sdfsprotocol.SNetPeerAddProtocol(&protocolBuilder, int32(protocolType))
	protocolBuilder.Finish(sdfsprotocol.SNetPeerEnd(&protocolBuilder))
	req.Param = protocolBuilder.Bytes[protocolBuilder.Head():]

	err = p.SNetClientDriver.Call(p.nameNodePeerID,
		"/DataNode/Register", &req, &resp)
	if err != nil {
		return err
	}

	var body = make([]byte, resp.BodySize)[:resp.BodySize]
	err = p.SNetClientDriver.ReadResponse(p.nameNodePeerID, &req, &resp, body)
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
