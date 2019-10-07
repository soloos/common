package solomqapi

import (
	"soloos/common/snet"
	"soloos/common/solofsapitypes"
	"soloos/common/solomqapitypes"
	"soloos/common/solomqprotocol"
	"soloos/common/soloosbase"
	"soloos/common/util"
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
	path string, ret interface{}, reqArgs ...interface{}) error {
	return p.SimpleCall(solodnPeerID,
		path, ret, reqArgs...)
}

func (p *SolomqClient) PrepareTopicNetBlockMetaData(peerID snet.PeerID,
	uTopic solomqapitypes.TopicUintptr,
	uNetBlock solofsapitypes.NetBlockUintptr,
	uNetINode solofsapitypes.NetINodeUintptr, netblockIndex int32) error {
	return nil
}

func (p *SolomqClient) UploadMemBlockWithSolomq(uTopic solomqapitypes.TopicUintptr,
	uJob solofsapitypes.UploadMemBlockJobUintptr,
	uploadPeerIndex int) error {
	var (
		snetReq            snet.SNetReq
		snetResp           snet.SNetResp
		req                solomqprotocol.TopicPWriteReq
		transferPeersCount int
		memBlockCap        int
		uploadChunkMask    offheap.ChunkMask
		respParamBs        []byte
		i                  int
		backendPeer        snet.Peer
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

		snetReq.Param = snet.MustSpecMarshalRequest(req)
		err = p.SrpcClientDriver.Call(backendPeer.ID,
			"/Topic/PWrite", &snetReq, &snetResp)
		if err != nil {
			goto QUERY_DONE
		}

		util.ChangeBytesArraySize(&respParamBs, int(snetResp.ParamSize))
		err = p.SrpcClientDriver.ReadResponse(backendPeer.ID, &snetReq, &snetResp, respParamBs, nil)
		if err != nil {
			goto QUERY_DONE
		}
	}

QUERY_DONE:
	return err
}
