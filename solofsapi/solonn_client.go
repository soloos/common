package solofsapi

import (
	"soloos/common/snet"
	"soloos/common/soloosbase"
)

type SolonnClient struct {
	*soloosbase.SoloosEnv
	solonnPeerID snet.PeerID
}

func (p *SolonnClient) Init(soloosEnv *soloosbase.SoloosEnv,
	solonnPeerID snet.PeerID) error {
	p.SoloosEnv = soloosEnv
	p.solonnPeerID = solonnPeerID
	return nil
}

func (p *SolonnClient) Dispatch(
	path string, resp snet.IResponse, reqArgs ...interface{}) error {
	// path string, resp snet.*Response, reqArgs ...interface{}) error {
	return p.SimpleCall(p.solonnPeerID,
		path, resp, reqArgs...)
}
