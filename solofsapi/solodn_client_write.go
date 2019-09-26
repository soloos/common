package solofsapi

import (
	"soloos/common/solofsapitypes"
	"soloos/common/solofsprotocol"
	"soloos/common/snettypes"
	"soloos/solodb/offheap"

	flatbuffers "github.com/google/flatbuffers/go"
)

func (p *SolodnClient) UploadMemBlock(uJob solofsapitypes.UploadMemBlockJobUintptr,
	uploadPeerIndex int,
) error {
	var solodn, _ = p.SoloosEnv.SNetDriver.GetPeer(
		uJob.Ptr().UNetBlock.Ptr().SyncDataBackends.Arr[uploadPeerIndex].PeerID)
	switch solodn.ServiceProtocol {
	case snettypes.ProtocolLocalFS:
		return p.uploadMemBlockWithDisk(uJob, uploadPeerIndex)
	case snettypes.ProtocolSolomq:
		return p.uploadMemBlockWithSolomq(uJob, uploadPeerIndex)
	case snettypes.ProtocolSolofs:
		return p.doUploadMemBlockWithSolofs(uJob, uploadPeerIndex)
	}

	return nil
}

func (p *SolodnClient) doUploadMemBlockWithSolofs(uJob solofsapitypes.UploadMemBlockJobUintptr,
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
		commonResp          solofsprotocol.CommonResponse
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
			backendOffs = backendOffs[:0]
			for i = 0; i < transferPeersCount; i++ {
				backendPeer, _ = p.SNetDriver.GetPeer(pNetBlock.SyncDataBackends.Arr[uploadPeerIndex+1+i].PeerID)
				backendOffs = append(backendOffs, protocolBuilder.CreateString(backendPeer.PeerIDStr()))
			}

			solofsprotocol.NetINodePWriteRequestStartTransferBackendsVector(&protocolBuilder, len(backendOffs))
			for i = len(backendOffs) - 1; i >= 0; i-- {
				protocolBuilder.PrependUOffsetT(backendOffs[i])
			}
			backendOff = protocolBuilder.EndVector(len(backendOffs))
		}

		netINodeIDOff = protocolBuilder.CreateByteVector(pNetBlock.NetINodeID[:])
		solofsprotocol.NetINodePWriteRequestStart(&protocolBuilder)
		if transferPeersCount > 0 {
			solofsprotocol.NetINodePWriteRequestAddTransferBackends(&protocolBuilder, backendOff)
		}
		solofsprotocol.NetINodePWriteRequestAddNetINodeID(&protocolBuilder, netINodeIDOff)
		solofsprotocol.NetINodePWriteRequestAddOffset(&protocolBuilder, uint64(netINodeWriteOffset))
		solofsprotocol.NetINodePWriteRequestAddLength(&protocolBuilder, int32(netINodeWriteLength))
		protocolBuilder.Finish(solofsprotocol.NetINodePWriteRequestEnd(&protocolBuilder))
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
			err = solofsapitypes.ErrNetBlockPWrite
			goto PWRITE_DONE
		}
		protocolBuilder.Reset()
	}

PWRITE_DONE:
	return err
}
