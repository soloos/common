package solofsapi

import (
	"soloos/common/solofsprotocol"
	"soloos/common/snettypes"

	flatbuffers "github.com/google/flatbuffers/go"
)

func SetNetINodeNetBlockInfoResponseError(protocolBuilder *flatbuffers.Builder, code int, err string) {
	protocolBuilder.Reset()
	var (
		errOff            flatbuffers.UOffsetT
		commonResponseOff flatbuffers.UOffsetT
	)
	errOff = protocolBuilder.CreateString(err)
	solofsprotocol.CommonResponseStart(protocolBuilder)
	solofsprotocol.CommonResponseAddCode(protocolBuilder, int32(code))
	solofsprotocol.CommonResponseAddError(protocolBuilder, errOff)
	commonResponseOff = solofsprotocol.CommonResponseEnd(protocolBuilder)

	solofsprotocol.NetINodeNetBlockInfoResponseStart(protocolBuilder)
	solofsprotocol.NetINodeNetBlockInfoResponseAddCommonResponse(protocolBuilder, commonResponseOff)
	protocolBuilder.Finish(solofsprotocol.NetINodeNetBlockInfoResponseEnd(protocolBuilder))
}

func SetNetINodeNetBlockInfoResponse(protocolBuilder *flatbuffers.Builder,
	backends []snettypes.PeerID, netBlockLen, netBlockCap int32) {
	var (
		backendOff        flatbuffers.UOffsetT
		commonResponseOff flatbuffers.UOffsetT
		i                 int
	)

	backendOffs := make([]flatbuffers.UOffsetT, len(backends))

	solofsprotocol.CommonResponseStart(protocolBuilder)
	solofsprotocol.CommonResponseAddCode(protocolBuilder, snettypes.CODE_OK)
	commonResponseOff = solofsprotocol.CommonResponseEnd(protocolBuilder)

	for i = 0; i < len(backends); i++ {
		backendOffs[i] = protocolBuilder.CreateString(backends[i].Str())
	}

	solofsprotocol.NetINodeNetBlockInfoResponseStartBackendsVector(protocolBuilder, len(backends))
	for i = len(backends) - 1; i >= 0; i-- {
		protocolBuilder.PrependUOffsetT(backendOffs[i])
	}
	backendOff = protocolBuilder.EndVector(len(backends))

	solofsprotocol.NetINodeNetBlockInfoResponseStart(protocolBuilder)
	solofsprotocol.NetINodeNetBlockInfoResponseAddCommonResponse(protocolBuilder, commonResponseOff)
	solofsprotocol.NetINodeNetBlockInfoResponseAddBackends(protocolBuilder, backendOff)
	solofsprotocol.NetINodeNetBlockInfoResponseAddLen(protocolBuilder, netBlockLen)
	solofsprotocol.NetINodeNetBlockInfoResponseAddCap(protocolBuilder, netBlockCap)
	protocolBuilder.Finish(solofsprotocol.NetINodeNetBlockInfoResponseEnd(protocolBuilder))
}
