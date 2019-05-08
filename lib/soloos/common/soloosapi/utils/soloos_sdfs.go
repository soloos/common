package base

func (p *SoloOS) initSDFS() error {
	var err error
	err = p.SDFSClientDriver.Init(&p.SoloOSEnv,
		p.options.SDFSNameNodeServeAddr,
		p.options.SDFSDBDriver, p.options.SDFSDsn,
	)
	if err != nil {
		return err
	}

	return nil
}
