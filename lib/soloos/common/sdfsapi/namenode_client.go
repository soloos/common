package sdfsapi

import (
	"soloos/common/snettypes"
	"soloos/common/soloosbase"
)

type NameNodeClient struct {
	*soloosbase.SoloOSEnv
	nameNodePeer snettypes.PeerUintptr
}

func (p *NameNodeClient) Init(soloOSEnv *soloosbase.SoloOSEnv,
	nameNodePeer snettypes.PeerUintptr) error {
	p.SoloOSEnv = soloOSEnv
	p.nameNodePeer = nameNodePeer
	return nil
}
