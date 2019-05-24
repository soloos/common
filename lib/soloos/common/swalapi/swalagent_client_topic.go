package swalapi

import (
	"soloos/common/sdfsapitypes"
	"soloos/common/snettypes"
	"soloos/common/swalapitypes"
)

func (p *SWALAgentClient) PrepareTopicFsINodeMetaData(
	pTopic *swalapitypes.Topic,
	pFsINodeMeta *sdfsapitypes.FsINodeMeta,
) error {
	// var (
	// req             snettypes.Request
	// resp            snettypes.Response
	// protocolBuilder flatbuffers.Builder
	// // uNetBlock sdfsapitypes.NetBlockUintptr
	// uPeer snettypes.PeerUintptr
	// )

	// err = p.SNetClientDriver.Call(uPeer,
	// "/Topic/PrepareFsINode", &req, &resp)

	return nil
}

func (p *SWALAgentClient) PrepareTopicNetBlockMetaData(peerID snettypes.PeerID,
	uTopic swalapitypes.TopicUintptr,
	uNetBlock sdfsapitypes.NetBlockUintptr,
	uNetINode sdfsapitypes.NetINodeUintptr, netblockIndex int32) error {
	return nil
}
