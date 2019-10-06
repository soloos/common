package solofsapi

import (
	"soloos/common/snettypes"
	"soloos/common/soloosbase"
)

type SolonnClient struct {
	*soloosbase.SoloosEnv
	solonnPeerID snettypes.PeerID
}

func (p *SolonnClient) Init(soloosEnv *soloosbase.SoloosEnv,
	solonnPeerID snettypes.PeerID) error {
	p.SoloosEnv = soloosEnv
	p.solonnPeerID = solonnPeerID
	return nil
}

func (p *SolonnClient) Dispatch(path string) error {
	return nil
}
