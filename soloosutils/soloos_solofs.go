package soloosutils

import (
	"soloos/common/log"
	"soloos/common/solofsapi"
	"soloos/common/snettypes"
)

func (p *SoloOS) initSOLOFS(solofsClientDriver solofsapi.ClientDriver) error {
	var err error
	var solonnPeerID snettypes.PeerID
	solonnPeerID.SetStr(p.options.SOLOFSSolonnSRPCPeerID)
	p.SOLOFSClientDriver = solofsClientDriver
	err = p.SOLOFSClientDriver.Init(&p.SoloOSEnv,
		solonnPeerID,
		p.options.SOLOFSDBDriver, p.options.SOLOFSDsn)
	if err != nil {
		log.Warn("SoloOS SOLOFSClientDriver Init error", err)
		return err
	}

	return nil
}
