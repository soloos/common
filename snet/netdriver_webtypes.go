package snet

import "soloos/common/snettypes"

type RegisterPeerReqJSON struct {
	PeerID   string `json:"PeerID"`
	Addr     string `json:"Addr"`
	Protocol string `json:"Protocol"`
}

type RegisterPeerRespJSON struct {
	snettypes.APIRespCommonJSON
}

type GetPeerReqJSON struct {
	PeerID string `json:"PeerID"`
}

type GetPeerRespJSON struct {
	snettypes.APIRespCommonJSON
	Data snettypes.PeerJSON `json:"Data"`
}
