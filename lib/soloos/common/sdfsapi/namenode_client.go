package sdfsapi

import (
	"soloos/common/snettypes"
	"soloos/common/soloosbase"
)

type NameNodeClient struct {
	*soloosbase.SoloOSEnv
	nameNodePeerID snettypes.PeerID
}

func (p *NameNodeClient) Init(soloOSEnv *soloosbase.SoloOSEnv,
	nameNodePeerID snettypes.PeerID) error {
	p.SoloOSEnv = soloOSEnv
	p.nameNodePeerID = nameNodePeerID
	return nil
}
