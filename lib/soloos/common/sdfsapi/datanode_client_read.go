package sdfsapi

import (
	"soloos/common/sdfsapitypes"
	"soloos/common/sdfsprotocol"
	"soloos/common/snettypes"

	flatbuffers "github.com/google/flatbuffers/go"
)

func (p *DataNodeClient) PReadMemBlock(uNetINode sdfsapitypes.NetINodeUintptr,
	uPeer snettypes.PeerUintptr,
	uNetBlock sdfsapitypes.NetBlockUintptr,
	netBlockIndex int32,
	uMemBlock sdfsapitypes.MemBlockUintptr,
	memBlockIndex int32,
	offset uint64, length int,
) (int, error) {
	if uNetBlock.Ptr().LocalDataBackend != 0 {
		return p.preadMemBlockWithDisk(uNetINode, uPeer, uNetBlock, netBlockIndex, uMemBlock, memBlockIndex, offset, length)
	}

	switch uPeer.Ptr().ServiceProtocol {
	case snettypes.ProtocolSRPC:
		return p.doPReadMemBlockWithSRPC(uNetINode, uPeer, uNetBlock, netBlockIndex, uMemBlock, memBlockIndex, offset, length)
	}

	return 0, sdfsapitypes.ErrServiceNotExists
}

func (p *DataNodeClient) doPReadMemBlockWithSRPC(uNetINode sdfsapitypes.NetINodeUintptr,
	uPeer snettypes.PeerUintptr,
	uNetBlock sdfsapitypes.NetBlockUintptr,
	netBlockIndex int32,
	uMemBlock sdfsapitypes.MemBlockUintptr,
	memBlockIndex int32,
	offset uint64, length int,
) (int, error) {
	var (
		req             snettypes.Request
		resp            snettypes.Response
		protocolBuilder flatbuffers.Builder
		netINodeIDOff   flatbuffers.UOffsetT
		err             error
	)

	netINodeIDOff = protocolBuilder.CreateByteVector(uNetBlock.Ptr().NetINodeID[:])
	sdfsprotocol.NetINodePReadRequestStart(&protocolBuilder)
	sdfsprotocol.NetINodePReadRequestAddNetINodeID(&protocolBuilder, netINodeIDOff)
	sdfsprotocol.NetINodePReadRequestAddOffset(&protocolBuilder, offset)
	sdfsprotocol.NetINodePReadRequestAddLength(&protocolBuilder, int32(length))
	protocolBuilder.Finish(sdfsprotocol.NetINodePReadRequestEnd(&protocolBuilder))
	req.Param = protocolBuilder.Bytes[protocolBuilder.Head():]

	// TODO choose datanode
	err = p.SNetClientDriver.Call(uPeer,
		"/NetINode/PRead", &req, &resp)
	if err != nil {
		return 0, err
	}

	var (
		netBlockPReadResp           sdfsprotocol.NetINodePReadResponse
		commonResp                  sdfsprotocol.CommonResponse
		param                       = make([]byte, resp.ParamSize)
		offsetInMemBlock, readedLen int
	)
	err = p.SNetClientDriver.ReadResponse(uPeer, &req, &resp, param)
	if err != nil {
		return 0, err
	}

	netBlockPReadResp.Init(param, flatbuffers.GetUOffsetT(param))
	netBlockPReadResp.CommonResponse(&commonResp)
	if commonResp.Code() != snettypes.CODE_OK {
		return 0, sdfsapitypes.ErrNetBlockPRead
	}

	offsetInMemBlock = int(offset - uint64(uMemBlock.Ptr().Bytes.Cap)*uint64(memBlockIndex))
	readedLen = int(resp.BodySize - resp.ParamSize)
	err = p.SNetClientDriver.ReadResponse(uPeer, &req, &resp,
		(*uMemBlock.Ptr().BytesSlice())[offsetInMemBlock:readedLen])
	if err != nil {
		return 0, err
	}

	return int(netBlockPReadResp.Length()), err
}
