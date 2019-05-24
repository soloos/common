package sdfsapi

import (
	"soloos/common/sdfsapitypes"
	"soloos/common/sdfsprotocol"
	"soloos/common/snettypes"
	"soloos/sdbone/offheap"

	flatbuffers "github.com/google/flatbuffers/go"
)

func (p *DataNodeClient) UploadMemBlock(uJob sdfsapitypes.UploadMemBlockJobUintptr,
	uploadPeerIndex int, transferPeersCount int,
) error {
	var (
		uDataNode snettypes.PeerUintptr
	)
	uDataNode = uJob.Ptr().UNetBlock.Ptr().SyncDataBackends.Arr[uploadPeerIndex]
	switch uDataNode.Ptr().ServiceProtocol {
	case snettypes.ProtocolDisk:
		return p.uploadMemBlockWithDisk(uJob, uploadPeerIndex, transferPeersCount)
	case snettypes.ProtocolSWAL:
		return p.uploadMemBlockWithSWAL(uJob, uploadPeerIndex, transferPeersCount)
	case snettypes.ProtocolSRPC:
		return p.doUploadMemBlockWithSRPC(uJob, uploadPeerIndex, transferPeersCount)
	}

	return nil
}

func (p *DataNodeClient) doUploadMemBlockWithSRPC(uJob sdfsapitypes.UploadMemBlockJobUintptr,
	uploadPeerIndex int, transferPeersCount int,
) error {
	var (
		req                 snettypes.Request
		resp                snettypes.Response
		protocolBuilder     flatbuffers.Builder
		netINodeIDOff       flatbuffers.UOffsetT
		backendOff          flatbuffers.UOffsetT
		uNetBlock           sdfsapitypes.NetBlockUintptr
		netINodeWriteOffset int
		netINodeWriteLength int
		memBlockCap         int
		peerOff, addrOff    flatbuffers.UOffsetT
		backendOffs         = make([]flatbuffers.UOffsetT, 8)
		uploadChunkMask     offheap.ChunkMask
		commonResp          sdfsprotocol.CommonResponse
		respBody            []byte
		i                   int
		uPeer               snettypes.PeerUintptr
		err                 error
	)

	uNetBlock = uJob.Ptr().UNetBlock
	uploadChunkMask = uJob.Ptr().GetProcessingChunkMask()

	req.OffheapBody.OffheapBytes = uJob.Ptr().UMemBlock.Ptr().Bytes.Data
	memBlockCap = uJob.Ptr().UMemBlock.Ptr().Bytes.Len
	for chunkMaskIndex := 0; chunkMaskIndex < uploadChunkMask.MaskArrayLen; chunkMaskIndex++ {
		req.OffheapBody.CopyOffset = uploadChunkMask.MaskArray[chunkMaskIndex].Offset
		req.OffheapBody.CopyEnd = uploadChunkMask.MaskArray[chunkMaskIndex].End
		netINodeWriteOffset = memBlockCap*int(uJob.Ptr().MemBlockIndex) + req.OffheapBody.CopyOffset
		netINodeWriteLength = req.OffheapBody.CopyEnd - req.OffheapBody.CopyOffset

		if transferPeersCount > 0 {
			for i = 0; i < transferPeersCount; i++ {
				uPeer = uNetBlock.Ptr().SyncDataBackends.Arr[uploadPeerIndex+1+i]
				peerOff = protocolBuilder.CreateByteVector(uPeer.Ptr().ID[:])
				addrOff = protocolBuilder.CreateString(uPeer.Ptr().AddressStr())
				sdfsprotocol.SNetPeerStart(&protocolBuilder)
				sdfsprotocol.SNetPeerAddPeerID(&protocolBuilder, peerOff)
				sdfsprotocol.SNetPeerAddAddress(&protocolBuilder, addrOff)
				if i < cap(backendOffs) {
					backendOffs[i] = sdfsprotocol.SNetPeerEnd(&protocolBuilder)
				} else {
					backendOffs = append(backendOffs, sdfsprotocol.SNetPeerEnd(&protocolBuilder))
				}
			}

			sdfsprotocol.NetINodePWriteRequestStartTransferBackendsVector(&protocolBuilder, transferPeersCount)
			for i = transferPeersCount - 1; i >= 0; i-- {
				protocolBuilder.PrependUOffsetT(backendOffs[i])
			}
			backendOff = protocolBuilder.EndVector(transferPeersCount)
		}

		netINodeIDOff = protocolBuilder.CreateByteVector(uNetBlock.Ptr().NetINodeID[:])
		sdfsprotocol.NetINodePWriteRequestStart(&protocolBuilder)
		if transferPeersCount > 0 {
			sdfsprotocol.NetINodePWriteRequestAddTransferBackends(&protocolBuilder, backendOff)
		}
		sdfsprotocol.NetINodePWriteRequestAddNetINodeID(&protocolBuilder, netINodeIDOff)
		sdfsprotocol.NetINodePWriteRequestAddOffset(&protocolBuilder, uint64(netINodeWriteOffset))
		sdfsprotocol.NetINodePWriteRequestAddLength(&protocolBuilder, int32(netINodeWriteLength))
		protocolBuilder.Finish(sdfsprotocol.NetINodePWriteRequestEnd(&protocolBuilder))
		req.Param = protocolBuilder.Bytes[protocolBuilder.Head():]

		uPeer = uJob.Ptr().UNetBlock.Ptr().SyncDataBackends.Arr[uploadPeerIndex]
		err = p.SNetClientDriver.Call(uPeer,
			"/NetINode/PWrite", &req, &resp)
		if err != nil {
			goto PWRITE_DONE
		}

		respBody = make([]byte, resp.ParamSize)
		err = p.SNetClientDriver.ReadResponse(uPeer, &req, &resp, respBody)
		if err != nil {
			goto PWRITE_DONE
		}
		commonResp.Init(respBody[:(resp.ParamSize)], flatbuffers.GetUOffsetT(respBody[:(resp.ParamSize)]))
		if commonResp.Code() != snettypes.CODE_OK {
			err = sdfsapitypes.ErrNetBlockPWrite
			goto PWRITE_DONE
		}
		protocolBuilder.Reset()
	}

PWRITE_DONE:
	return err
}
