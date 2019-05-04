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

	err = p.SDFSClientDriver.InitClient(&p.SDFSClient,
		defautlNetBlockCap,
		defautlMemBlockCap, defautlMemBlocksLimit)
	if err != nil {
		return err
	}

	p.PosixFS = p.SDFSClient.GetPosixFS()

	return nil
}
