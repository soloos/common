package swalapi

import "soloos/common/soloosbase"

type SWALAgentClient struct {
	*soloosbase.SoloOSEnv
}

func (p *SWALAgentClient) Init(soloOSEnv *soloosbase.SoloOSEnv) error {
	p.SoloOSEnv = soloOSEnv
	return nil
}
