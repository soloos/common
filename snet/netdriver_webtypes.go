package snet

import "soloos/common/snettypes"

type RegisterPeerReq struct {
	PeerID   string
	Addr     string
	Protocol string
}

type RegisterPeerResp struct {
	snettypes.RespDataCommon
}

type GetPeerResp struct {
	snettypes.RespDataCommon
	Data snettypes.PeerJSON
}
