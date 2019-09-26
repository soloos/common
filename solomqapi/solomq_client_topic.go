package solomqapi

import (
	"soloos/common/log"
	"soloos/common/solofsapitypes"
	"soloos/common/solofsprotocol"
	"soloos/common/snettypes"
	"soloos/common/solomqapitypes"
	"soloos/common/solomqprotocol"
	"soloos/solodb/offheap"
	"time"

	flatbuffers "github.com/google/flatbuffers/go"
)

func (p *SolomqClient) PrepareTopicMetaData(peerID snettypes.PeerID,
	uTopic solomqapitypes.TopicUintptr,
	fsINodeID solofsapitypes.FsINodeID,
) error {
	var (
		req             snettypes.Request
		resp            snettypes.Response
		protocolBuilder flatbuffers.Builder
		pTopic          = uTopic.Ptr()
		commonResp      solomqprotocol.CommonResponse
		respBody        []byte
		err             error
	)

	solomqprotocol.TopicPrepareRequestStart(&protocolBuilder)
	solomqprotocol.TopicPrepareRequestAddTopicID(&protocolBuilder, pTopic.ID)
	solomqprotocol.TopicPrepareRequestAddFsINodeID(&protocolBuilder, fsINodeID)
	protocolBuilder.Finish(solomqprotocol.TopicPrepareRequestEnd(&protocolBuilder))
	req.Param = protocolBuilder.Bytes[protocolBuilder.Head():]

	for i := 0; i < p.normalCallRetryTimes; i++ {
		err = p.SNetClientDriver.Call(peerID,
			"/Topic/Prepare", &req, &resp)
		if err == nil {
			break
		}
		log.Info("Topic/Prepare peerID:", peerID.Str(),
			", topicid:", pTopic.ID,
			", retryTimes:", i,
			", err", err)
		time.Sleep(p.waitAliveEveryRetryWaitTime)
	}
	if err != nil {
		return err
	}

	respBody = make([]byte, resp.ParamSize)
	commonResp.Init(respBody[:(resp.ParamSize)], flatbuffers.GetUOffsetT(respBody[:(resp.ParamSize)]))
	if commonResp.Code() != snettypes.CODE_OK {
		err = solofsapitypes.ErrNetBlockPWrite
		goto QUERY_DONE
	}
	protocolBuilder.Reset()

QUERY_DONE:
	return err
}

func (p *SolomqClient) PrepareTopicNetBlockMetaData(peerID snettypes.PeerID,
	uTopic solomqapitypes.TopicUintptr,
	uNetBlock solofsapitypes.NetBlockUintptr,
	uNetINode solofsapitypes.NetINodeUintptr, netblockIndex int32) error {
	return nil
}

func (p *SolomqClient) UploadMemBlockWithSolomq(uTopic solomqapitypes.TopicUintptr,
	uJob solofsapitypes.UploadMemBlockJobUintptr,
	uploadPeerIndex int) error {
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
		commonResp          solomqprotocol.CommonResponse
		respBody            []byte
		i                   int
		backendPeer         snettypes.Peer
		err                 error
	)

	var pJob = uJob.Ptr()
	var pMemBlock = pJob.UMemBlock.Ptr()
	var pNetBlock = uJob.Ptr().UNetBlock.Ptr()
	var pTopic = uTopic.Ptr()
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
		solomqprotocol.TopicPWriteRequestStart(&protocolBuilder)
		if transferPeersCount > 0 {
			solomqprotocol.TopicPWriteRequestAddTransferBackends(&protocolBuilder, backendOff)
		}
		solomqprotocol.TopicPWriteRequestAddTopicID(&protocolBuilder, pTopic.ID)
		solomqprotocol.TopicPWriteRequestAddNetINodeID(&protocolBuilder, netINodeIDOff)
		solomqprotocol.TopicPWriteRequestAddOffset(&protocolBuilder, uint64(netINodeWriteOffset))
		solomqprotocol.TopicPWriteRequestAddLength(&protocolBuilder, int32(netINodeWriteLength))
		protocolBuilder.Finish(solomqprotocol.TopicPWriteRequestEnd(&protocolBuilder))
		req.Param = protocolBuilder.Bytes[protocolBuilder.Head():]

		backendPeer, err = p.SNetDriver.GetPeer(pNetBlock.SyncDataBackends.Arr[uploadPeerIndex].PeerID)
		if err != nil {
			goto QUERY_DONE
		}

		err = p.SNetClientDriver.Call(backendPeer.ID,
			"/Topic/PWrite", &req, &resp)
		if err != nil {
			goto QUERY_DONE
		}

		respBody = make([]byte, resp.ParamSize)
		err = p.SNetClientDriver.ReadResponse(backendPeer.ID, &req, &resp, respBody)
		if err != nil {
			goto QUERY_DONE
		}
		commonResp.Init(respBody[:(resp.ParamSize)], flatbuffers.GetUOffsetT(respBody[:(resp.ParamSize)]))
		if commonResp.Code() != snettypes.CODE_OK {
			err = solofsapitypes.ErrNetBlockPWrite
			goto QUERY_DONE
		}
		protocolBuilder.Reset()
	}

QUERY_DONE:
	return err
}
