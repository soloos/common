package sdfsapi

import (
	"soloos/common/snettypes"
	"soloos/common/sdfsprotocol"

	flatbuffers "github.com/google/flatbuffers/go"
)

func (p *NameNodeClient) RegisterDataNode(peerID snettypes.PeerID, serveAddr string) error {
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
	protocolBuilder.Finish(sdfsprotocol.SNetPeerEnd(&protocolBuilder))
	req.Param = protocolBuilder.Bytes[protocolBuilder.Head():]

	err = p.SNetClientDriver.Call(p.nameNodePeer,
		"/DataNode/Register", &req, &resp)
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
