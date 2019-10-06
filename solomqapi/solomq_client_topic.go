package solomqapi

import (
	"soloos/common/log"
	"soloos/common/snettypes"
	"soloos/common/solofsapitypes"
	"soloos/common/solomqapitypes"
	"soloos/common/solomqprotocol"
	"soloos/common/util"
	"soloos/solodb/offheap"
	"time"
)

func (p *SolomqClient) PrepareTopicMetaData(peerID snettypes.PeerID,
	uTopic solomqapitypes.TopicUintptr,
	fsINodeID solofsapitypes.FsINodeID,
) error {
	var (
		req    solomqprotocol.TopicPrepareReq
		pTopic = uTopic.Ptr()
		err    error
	)

	req.TopicID = pTopic.ID
	req.FsINodeID = fsINodeID

	for i := 0; i < p.normalCallRetryTimes; i++ {
		err = p.SNetClientDriver.SimpleCall(peerID,
			"/Topic/Prepare", nil, req)
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

	return nil
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
		snetReq            snettypes.SNetReq
		snetResp           snettypes.SNetResp
		req                solomqprotocol.TopicPWriteReq
		transferPeersCount int
		memBlockCap        int
		uploadChunkMask    offheap.ChunkMask
		respParamBs        []byte
		i                  int
		backendPeer        snettypes.Peer
		err                error
	)

	var pJob = uJob.Ptr()
	var pMemBlock = pJob.UMemBlock.Ptr()
	var pNetBlock = uJob.Ptr().UNetBlock.Ptr()
	var pTopic = uTopic.Ptr()
	uploadChunkMask = pJob.GetProcessingChunkMask()
	transferPeersCount = int(pNetBlock.SyncDataBackends.Arr[uploadPeerIndex].TransferCount)

	snetReq.OffheapBody.OffheapBytes = pMemBlock.Bytes.Data
	memBlockCap = pMemBlock.Bytes.Len
	for chunkMaskIndex := 0; chunkMaskIndex < uploadChunkMask.MaskArrayLen; chunkMaskIndex++ {
		snetReq.OffheapBody.CopyOffset = uploadChunkMask.MaskArray[chunkMaskIndex].Offset
		snetReq.OffheapBody.CopyEnd = uploadChunkMask.MaskArray[chunkMaskIndex].End

		req.TopicID = pTopic.ID
		req.Offset = uint64(memBlockCap)*uint64(pJob.MemBlockIndex) + uint64(snetReq.OffheapBody.CopyOffset)
		req.Length = snetReq.OffheapBody.CopyEnd - snetReq.OffheapBody.CopyOffset
		req.TransferBackends = req.TransferBackends[:0]
		for i = 0; i < transferPeersCount; i++ {
			backendPeer, _ = p.SNetDriver.GetPeer(pNetBlock.SyncDataBackends.Arr[uploadPeerIndex+1+i].PeerID)
			req.TransferBackends = append(req.TransferBackends, backendPeer.PeerIDStr())
		}

		backendPeer, err = p.SNetDriver.GetPeer(pNetBlock.SyncDataBackends.Arr[uploadPeerIndex].PeerID)
		if err != nil {
			goto QUERY_DONE
		}

		err = p.SNetClientDriver.Call(backendPeer.ID,
			"/Topic/PWrite", &snetReq, &snetResp, req)
		if err != nil {
			goto QUERY_DONE
		}

		util.ChangeBytesArraySize(&respParamBs, int(snetResp.ParamSize))
		err = p.SNetClientDriver.ReadResponse(backendPeer.ID, &snetReq, &snetResp, respParamBs, nil)
		if err != nil {
			goto QUERY_DONE
		}
	}

QUERY_DONE:
	return err
}
