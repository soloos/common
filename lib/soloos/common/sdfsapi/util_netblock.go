package sdfsapi

import (
	"soloos/common/snettypes"
	"soloos/common/sdfsprotocol"

	flatbuffers "github.com/google/flatbuffers/go"
)

func SetNetINodeNetBlockInfoResponseError(protocolBuilder *flatbuffers.Builder, code int, err string) {
	protocolBuilder.Reset()
	var (
		errOff            flatbuffers.UOffsetT
		commonResponseOff flatbuffers.UOffsetT
	)
	errOff = protocolBuilder.CreateString(err)
	sdfsprotocol.CommonResponseStart(protocolBuilder)
	sdfsprotocol.CommonResponseAddCode(protocolBuilder, int32(code))
	sdfsprotocol.CommonResponseAddError(protocolBuilder, errOff)
	commonResponseOff = sdfsprotocol.CommonResponseEnd(protocolBuilder)

	sdfsprotocol.NetINodeNetBlockInfoResponseStart(protocolBuilder)
	sdfsprotocol.NetINodeNetBlockInfoResponseAddCommonResponse(protocolBuilder, commonResponseOff)
	protocolBuilder.Finish(sdfsprotocol.NetINodeNetBlockInfoResponseEnd(protocolBuilder))
}

func SetNetINodeNetBlockInfoResponse(protocolBuilder *flatbuffers.Builder,
	backends []snettypes.PeerUintptr, netBlockLen, netBlockCap int32) {
	var (
		peerOff           flatbuffers.UOffsetT
		addrOff           flatbuffers.UOffsetT
		backendOff        flatbuffers.UOffsetT
		commonResponseOff flatbuffers.UOffsetT
		i                 int
	)

	backendOffs := make([]flatbuffers.UOffsetT, len(backends))

	sdfsprotocol.CommonResponseStart(protocolBuilder)
	sdfsprotocol.CommonResponseAddCode(protocolBuilder, snettypes.CODE_OK)
	commonResponseOff = sdfsprotocol.CommonResponseEnd(protocolBuilder)

	for i = 0; i < len(backends); i++ {
		peerOff = protocolBuilder.CreateByteVector(backends[i].Ptr().ID[:])
		addrOff = protocolBuilder.CreateString(backends[i].Ptr().AddressStr())
		sdfsprotocol.SNetPeerStart(protocolBuilder)
		sdfsprotocol.SNetPeerAddPeerID(protocolBuilder, peerOff)
		sdfsprotocol.SNetPeerAddAddress(protocolBuilder, addrOff)
		backendOffs[i] = sdfsprotocol.SNetPeerEnd(protocolBuilder)
	}

	sdfsprotocol.NetINodeNetBlockInfoResponseStartBackendsVector(protocolBuilder, len(backends))
	for i = len(backends) - 1; i >= 0; i-- {
		protocolBuilder.PrependUOffsetT(backendOffs[i])
	}
	backendOff = protocolBuilder.EndVector(len(backends))

	sdfsprotocol.NetINodeNetBlockInfoResponseStart(protocolBuilder)
	sdfsprotocol.NetINodeNetBlockInfoResponseAddCommonResponse(protocolBuilder, commonResponseOff)
	sdfsprotocol.NetINodeNetBlockInfoResponseAddBackends(protocolBuilder, backendOff)
	sdfsprotocol.NetINodeNetBlockInfoResponseAddLen(protocolBuilder, netBlockLen)
	sdfsprotocol.NetINodeNetBlockInfoResponseAddCap(protocolBuilder, netBlockCap)
	protocolBuilder.Finish(sdfsprotocol.NetINodeNetBlockInfoResponseEnd(protocolBuilder))
}
