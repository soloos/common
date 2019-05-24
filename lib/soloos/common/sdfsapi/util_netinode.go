package sdfsapi

import (
	"soloos/common/sdfsprotocol"
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
	sdfsprotocol.CommonResponseStart(protocolBuilder)
	sdfsprotocol.CommonResponseAddCode(protocolBuilder, int32(code))
	sdfsprotocol.CommonResponseAddError(protocolBuilder, errOff)
	commonResponseOff = sdfsprotocol.CommonResponseEnd(protocolBuilder)

	sdfsprotocol.NetINodeInfoResponseStart(protocolBuilder)
	sdfsprotocol.NetINodeInfoResponseAddCommonResponse(protocolBuilder, commonResponseOff)
	protocolBuilder.Finish(sdfsprotocol.NetINodeInfoResponseEnd(protocolBuilder))
}

func SetNetINodeInfoResponse(protocolBuilder *flatbuffers.Builder,
	size uint64, netBlockCap int32, memBlockCap int32) {
	protocolBuilder.Reset()
	var (
		commonResponseOff flatbuffers.UOffsetT
	)
	sdfsprotocol.CommonResponseStart(protocolBuilder)
	sdfsprotocol.CommonResponseAddCode(protocolBuilder, snettypes.CODE_OK)
	commonResponseOff = sdfsprotocol.CommonResponseEnd(protocolBuilder)

	sdfsprotocol.NetINodeInfoResponseStart(protocolBuilder)
	sdfsprotocol.NetINodeInfoResponseAddCommonResponse(protocolBuilder, commonResponseOff)
	sdfsprotocol.NetINodeInfoResponseAddSize(protocolBuilder, size)
	sdfsprotocol.NetINodeInfoResponseAddNetBlockCap(protocolBuilder, int32(netBlockCap))
	sdfsprotocol.NetINodeInfoResponseAddMemBlockCap(protocolBuilder, int32(memBlockCap))
	protocolBuilder.Finish(sdfsprotocol.NetINodeInfoResponseEnd(protocolBuilder))
}

func SetNetINodePReadResponseError(protocolBuilder *flatbuffers.Builder, code int, err string) {
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

	sdfsprotocol.NetINodePReadResponseStart(protocolBuilder)
	sdfsprotocol.NetINodePReadResponseAddCommonResponse(protocolBuilder, commonResponseOff)
	protocolBuilder.Finish(sdfsprotocol.NetINodePReadResponseEnd(protocolBuilder))
}

func SetNetINodePReadResponse(protocolBuilder *flatbuffers.Builder, length int32) {
	protocolBuilder.Reset()
	var (
		commonResponseOff flatbuffers.UOffsetT
	)
	sdfsprotocol.CommonResponseStart(protocolBuilder)
	sdfsprotocol.CommonResponseAddCode(protocolBuilder, snettypes.CODE_OK)
	commonResponseOff = sdfsprotocol.CommonResponseEnd(protocolBuilder)

	sdfsprotocol.NetINodePReadResponseStart(protocolBuilder)
	sdfsprotocol.NetINodePReadResponseAddCommonResponse(protocolBuilder, commonResponseOff)
	sdfsprotocol.NetINodePReadResponseAddLength(protocolBuilder, length)
	protocolBuilder.Finish(sdfsprotocol.NetINodePReadResponseEnd(protocolBuilder))
}
