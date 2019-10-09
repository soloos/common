package solomqapi

import (
	"soloos/common/snet"
	"soloos/common/solofstypes"
	"soloos/common/solomqprotocol"
	"soloos/common/solomqtypes"
	"soloos/common/soloosbase"
	"soloos/solodb/offheap"
)

type SolomqClient struct {
	*soloosbase.SoloosEnv
}

func (p *SolomqClient) Init(soloosEnv *soloosbase.SoloosEnv) error {
	p.SoloosEnv = soloosEnv
	return nil
}

func (p *SolomqClient) Dispatch(solodnPeerID snet.PeerID,
	path string, resp *snet.Response, reqArgs ...interface{}) error {
	return p.SimpleCall(solodnPeerID,
		path, resp, reqArgs...)
}

func (p *SolomqClient) PrepareTopicNetBlockMetaData(peerID snet.PeerID,
	uTopic solomqtypes.TopicUintptr,
	uNetBlock solofstypes.NetBlockUintptr,
	uNetINode solofstypes.NetINodeUintptr, netblockIndex int32) error {
	return nil
}

func (p *SolomqClient) UploadMemBlockWithSolomq(uTopic solomqtypes.TopicUintptr,
	uJob solofstypes.UploadMemBlockJobUintptr,
	uploadPeerIndex int) error {
	var (
		snetReq            snet.SNetReq
		snetResp           snet.SNetResp
		req                solomqprotocol.TopicPWriteReq
		transferPeersCount int
		memBlockCap        int
		uploadChunkMask    offheap.ChunkMask
		i                  int
		backendPeer        snet.Peer
		err                error
	)

	var pJob = uJob.Ptr()
	var pMemBlock = pJob.UMemBlock.Ptr()
	var pNetBlock = uJob.Ptr().UNetBlock.Ptr()
	var pTopic = uTopic.Ptr()
	var resp = snet.Response{RespData: ""}
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

		err = p.SrpcClientDriver.Call(backendPeer.ID,
			"/Topic/PWrite", &snetReq, &snetResp, req)
		if err != nil {
			goto QUERY_DONE
		}

		err = p.SrpcClientDriver.SimpleReadResponse(backendPeer.ID, &snetReq, &snetResp, &resp)
		if err != nil {
			goto QUERY_DONE
		}
	}

QUERY_DONE:
	return err
}
