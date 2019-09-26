package solofsapi

import (
	"soloos/common/snettypes"
	"soloos/common/soloosbase"
)

type SolonnClient struct {
	*soloosbase.SoloOSEnv
	solonnPeerID snettypes.PeerID
}

func (p *SolonnClient) Init(soloOSEnv *soloosbase.SoloOSEnv,
	solonnPeerID snettypes.PeerID) error {
	p.SoloOSEnv = soloOSEnv
	p.solonnPeerID = solonnPeerID
	return nil
}
