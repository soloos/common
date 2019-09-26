package solofsapi

import (
	"soloos/common/solofsprotocol"
	"soloos/common/snettypes"

	flatbuffers "github.com/google/flatbuffers/go"
)

func (p *SolonnClient) SolodnRegister(peerID snettypes.PeerID,
	serveAddr string,
	protocolType snettypes.ServiceProtocol) error {
	var (
		req             snettypes.Request
		resp            snettypes.Response
		protocolBuilder flatbuffers.Builder
		peerIDOff       flatbuffers.UOffsetT
		addrOff         flatbuffers.UOffsetT
		protocolOff     flatbuffers.UOffsetT
		err             error
	)

	peerIDOff = protocolBuilder.CreateByteString(peerID[:])
	addrOff = protocolBuilder.CreateString(serveAddr)
	protocolOff = protocolBuilder.CreateString(protocolType.Str())
	solofsprotocol.SNetPeerStart(&protocolBuilder)
	solofsprotocol.SNetPeerAddPeerID(&protocolBuilder, peerIDOff)
	solofsprotocol.SNetPeerAddAddress(&protocolBuilder, addrOff)
	solofsprotocol.SNetPeerAddProtocol(&protocolBuilder, protocolOff)
	protocolBuilder.Finish(solofsprotocol.SNetPeerEnd(&protocolBuilder))
	req.Param = protocolBuilder.Bytes[protocolBuilder.Head():]

	err = p.SNetClientDriver.Call(p.solonnPeerID,
		"/Solodn/Register", &req, &resp)
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
