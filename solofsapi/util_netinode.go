package solofsapi

import (
	"soloos/common/solofsprotocol"
	"soloos/common/snettypes"

	flatbuffers "github.com/google/flatbuffers/go"
)

func SetNetINodeInfoResponseError(protocolBuilder *flatbuffers.Builder, code int, err string) {
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

	solofsprotocol.NetINodeInfoResponseStart(protocolBuilder)
	solofsprotocol.NetINodeInfoResponseAddCommonResponse(protocolBuilder, commonResponseOff)
	protocolBuilder.Finish(solofsprotocol.NetINodeInfoResponseEnd(protocolBuilder))
}

func SetNetINodeInfoResponse(protocolBuilder *flatbuffers.Builder,
	size uint64, netBlockCap int32, memBlockCap int32) {
	protocolBuilder.Reset()
	var (
		commonResponseOff flatbuffers.UOffsetT
	)
	solofsprotocol.CommonResponseStart(protocolBuilder)
	solofsprotocol.CommonResponseAddCode(protocolBuilder, snettypes.CODE_OK)
	commonResponseOff = solofsprotocol.CommonResponseEnd(protocolBuilder)

	solofsprotocol.NetINodeInfoResponseStart(protocolBuilder)
	solofsprotocol.NetINodeInfoResponseAddCommonResponse(protocolBuilder, commonResponseOff)
	solofsprotocol.NetINodeInfoResponseAddSize(protocolBuilder, size)
	solofsprotocol.NetINodeInfoResponseAddNetBlockCap(protocolBuilder, int32(netBlockCap))
	solofsprotocol.NetINodeInfoResponseAddMemBlockCap(protocolBuilder, int32(memBlockCap))
	protocolBuilder.Finish(solofsprotocol.NetINodeInfoResponseEnd(protocolBuilder))
}

func SetNetINodePReadResponseError(protocolBuilder *flatbuffers.Builder, code int, err string) {
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

	solofsprotocol.NetINodePReadResponseStart(protocolBuilder)
	solofsprotocol.NetINodePReadResponseAddCommonResponse(protocolBuilder, commonResponseOff)
	protocolBuilder.Finish(solofsprotocol.NetINodePReadResponseEnd(protocolBuilder))
}

func SetNetINodePReadResponse(protocolBuilder *flatbuffers.Builder, length int32) {
	protocolBuilder.Reset()
	var (
		commonResponseOff flatbuffers.UOffsetT
	)
	solofsprotocol.CommonResponseStart(protocolBuilder)
	solofsprotocol.CommonResponseAddCode(protocolBuilder, snettypes.CODE_OK)
	commonResponseOff = solofsprotocol.CommonResponseEnd(protocolBuilder)

	solofsprotocol.NetINodePReadResponseStart(protocolBuilder)
	solofsprotocol.NetINodePReadResponseAddCommonResponse(protocolBuilder, commonResponseOff)
	solofsprotocol.NetINodePReadResponseAddLength(protocolBuilder, length)
	protocolBuilder.Finish(solofsprotocol.NetINodePReadResponseEnd(protocolBuilder))
}
