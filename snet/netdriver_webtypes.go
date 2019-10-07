package snet

type RegisterPeerReq struct {
	PeerID   string
	Addr     string
	Protocol string
}

type RegisterPeerResp struct {
	RespDataCommon
}

type GetPeerResp struct {
	RespDataCommon
	Data PeerJSON
}
