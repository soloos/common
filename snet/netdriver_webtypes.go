package snet

type RegisterPeerReq struct {
	PeerID   string
	Addr     string
	Protocol string
}

type RegisterPeerResp struct {
	RespCommon
}

type GetPeerResp struct {
	RespCommon
	Data PeerJSON
}
