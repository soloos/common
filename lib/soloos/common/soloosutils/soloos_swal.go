package soloosutils

func (p *SoloOS) initSWAL() error {
	var err error
	err = p.SWALClientDriver.Init(&p.SoloOSEnv,
		p.options.SWALAgentPeerID, p.options.SWALAgentServeAddr,
		p.options.SWALDBDriver, p.options.SWALDsn,
		p.options.SWALDefaultNetBlockCap, p.options.SWALDefaultMemBlockCap,
	)
	if err != nil {
		return err
	}

	return nil
}
