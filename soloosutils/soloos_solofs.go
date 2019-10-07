package soloosutils

import (
	"soloos/common/log"
	"soloos/common/snet"
	"soloos/common/solofsapi"
)

func (p *Soloos) initSolofs(solofsClientDriver solofsapi.ClientDriver) error {
	var err error
	var solonnPeerID snet.PeerID
	solonnPeerID.SetStr(p.options.SolofsSolonnSrpcPeerID)
	p.SolofsClientDriver = solofsClientDriver
	err = p.SolofsClientDriver.Init(&p.SoloosEnv,
		solonnPeerID,
		p.options.SolofsDBDriver, p.options.SolofsDsn)
	if err != nil {
		log.Warn("Soloos SolofsClientDriver Init error", err)
		return err
	}

	return nil
}
