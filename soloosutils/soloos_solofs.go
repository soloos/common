package soloosutils

import (
	"soloos/common/log"
	"soloos/common/solofsapi"
	"soloos/common/snettypes"
)

func (p *Soloos) initSolofs(solofsClientDriver solofsapi.ClientDriver) error {
	var err error
	var solonnPeerID snettypes.PeerID
	solonnPeerID.SetStr(p.options.SolofsSolonnSRPCPeerID)
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
