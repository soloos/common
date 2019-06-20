package soloosutils

import "soloos/common/swalapi"

func (p *SoloOS) initSWAL(swalClientDriver swalapi.ClientDriver) error {
	var err error
	p.SWALClientDriver = swalClientDriver
	err = p.SWALClientDriver.Init(&p.SoloOSEnv,
		p.options.SWALAgentPeerID, p.options.SWALAgentServeAddr,
		p.options.SWALDBDriver, p.options.SWALDsn,
		p.options.SWALDefaultNetBlockCap, p.options.SWALDefaultMemBlockCap)
	if err != nil {
		return err
	}

	return nil
}
