package sdfsapi

import (
	"soloos/common/sdfsapitypes"
	"soloos/common/sdfsprotocol"
	"soloos/common/snettypes"
	"soloos/sdbone/offheap"

	flatbuffers "github.com/google/flatbuffers/go"
)

func (p *DataNodeClient) UploadMemBlock(uJob sdfsapitypes.UploadMemBlockJobUintptr,
	uploadPeerIndex int,
) error {
	var dataNode, _ = p.SoloOSEnv.SNetDriver.GetPeer(
		uJob.Ptr().UNetBlock.Ptr().SyncDataBackends.Arr[uploadPeerIndex].PeerID)
	switch dataNode.ServiceProtocol {
	case snettypes.ProtocolDisk:
		return p.uploadMemBlockWithDisk(uJob, uploadPeerIndex)
	case snettypes.ProtocolSWAL:
		return p.uploadMemBlockWithSWAL(uJob, uploadPeerIndex)
	case snettypes.ProtocolSDFS:
		return p.doUploadMemBlockWithSDFS(uJob, uploadPeerIndex)
	}

	return nil
}

func (p *DataNodeClient) doUploadMemBlockWithSDFS(uJob sdfsapitypes.UploadMemBlockJobUintptr,
	uploadPeerIndex int,
) error {
	var (
		req                 snettypes.Request
		resp                snettypes.Response
		protocolBuilder     flatbuffers.Builder
		netINodeIDOff       flatbuffers.UOffsetT
		backendOff          flatbuffers.UOffsetT
		transferPeersCount  int
		netINodeWriteOffset int
		netINodeWriteLength int
		memBlockCap         int
		backendOffs         = make([]flatbuffers.UOffsetT, 8)
		uploadChunkMask     offheap.ChunkMask
		commonResp          sdfsprotocol.CommonResponse
		respBody            []byte
		i                   int
		backendPeer         snettypes.Peer
		err                 error
	)

	var pJob = uJob.Ptr()
	var pNetBlock = pJob.UNetBlock.Ptr()
	var pMemBlock = pJob.UMemBlock.Ptr()
	uploadChunkMask = pJob.GetProcessingChunkMask()
	transferPeersCount = int(pNetBlock.SyncDataBackends.Arr[uploadPeerIndex].TransferCount)

	req.OffheapBody.OffheapBytes = pMemBlock.Bytes.Data
	memBlockCap = pMemBlock.Bytes.Len
	for chunkMaskIndex := 0; chunkMaskIndex < uploadChunkMask.MaskArrayLen; chunkMaskIndex++ {
		req.OffheapBody.CopyOffset = uploadChunkMask.MaskArray[chunkMaskIndex].Offset
		req.OffheapBody.CopyEnd = uploadChunkMask.MaskArray[chunkMaskIndex].End
		netINodeWriteOffset = memBlockCap*int(pJob.MemBlockIndex) + req.OffheapBody.CopyOffset
		netINodeWriteLength = req.OffheapBody.CopyEnd - req.OffheapBody.CopyOffset

		if transferPeersCount > 0 {
			for i = 0; i < transferPeersCount; i++ {
				backendPeer, _ = p.SNetDriver.GetPeer(pNetBlock.SyncDataBackends.Arr[uploadPeerIndex+1+i].PeerID)
				if i < cap(backendOffs) {
					backendOffs[i] = protocolBuilder.CreateString(backendPeer.PeerIDStr())
				} else {
					backendOffs = append(backendOffs, protocolBuilder.CreateString(backendPeer.PeerIDStr()))
				}
			}

			sdfsprotocol.NetINodePWriteRequestStartTransferBackendsVector(&protocolBuilder, len(backendOffs))
			for i = len(backendOffs) - 1; i >= 0; i-- {
				protocolBuilder.PrependUOffsetT(backendOffs[i])
			}
			backendOff = protocolBuilder.EndVector(len(backendOffs))
		}

		netINodeIDOff = protocolBuilder.CreateByteVector(pNetBlock.NetINodeID[:])
		sdfsprotocol.NetINodePWriteRequestStart(&protocolBuilder)
		if transferPeersCount > 0 {
			sdfsprotocol.NetINodePWriteRequestAddTransferBackends(&protocolBuilder, backendOff)
		}
		sdfsprotocol.NetINodePWriteRequestAddNetINodeID(&protocolBuilder, netINodeIDOff)
		sdfsprotocol.NetINodePWriteRequestAddOffset(&protocolBuilder, uint64(netINodeWriteOffset))
		sdfsprotocol.NetINodePWriteRequestAddLength(&protocolBuilder, int32(netINodeWriteLength))
		protocolBuilder.Finish(sdfsprotocol.NetINodePWriteRequestEnd(&protocolBuilder))
		req.Param = protocolBuilder.Bytes[protocolBuilder.Head():]

		backendPeer, err = p.SNetDriver.GetPeer(pNetBlock.SyncDataBackends.Arr[uploadPeerIndex].PeerID)
		if err != nil {
			goto PWRITE_DONE
		}

		err = p.SNetClientDriver.Call(backendPeer.ID,
			"/NetINode/PWrite", &req, &resp)
		if err != nil {
			goto PWRITE_DONE
		}

		respBody = make([]byte, resp.ParamSize)
		err = p.SNetClientDriver.ReadResponse(backendPeer.ID, &req, &resp, respBody)
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
