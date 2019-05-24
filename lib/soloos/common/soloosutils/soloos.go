package soloosutils

import (
	"soloos/common/sdfsapi"
	"soloos/common/soloosbase"
	"soloos/common/swalapi"
)

var (
	SoloOSInstance        SoloOS
	isDefaultSoloOSInited bool = false
)

type SoloOS struct {
	options Options
	soloosbase.SoloOSEnv

	SDFSClientDriver sdfsapi.ClientDriver
	SWALClientDriver swalapi.ClientDriver
}

func InitSoloOSInstance(options Options,
	sdfsClientDriver sdfsapi.ClientDriver,
	swalClientDriver swalapi.ClientDriver,
) error {
	if isDefaultSoloOSInited {
		return nil
	}
	isDefaultSoloOSInited = true
	return SoloOSInstance.Init(options, sdfsClientDriver, swalClientDriver)
}

func (p *SoloOS) Init(options Options,
	sdfsClientDriver sdfsapi.ClientDriver,
	swalClientDriver swalapi.ClientDriver,
) error {
	var err error

	p.options = options
	err = p.SoloOSEnv.Init()
	if err != nil {
		return err
	}

	p.SDFSClientDriver = sdfsClientDriver
	err = p.SDFSClientDriver.Init(&p.SoloOSEnv,
		options.SDFSNameNodeServeAddr,
		options.SDFSDBDriver, options.SDFSDsn)
	if err != nil {
		return err
	}

	p.SWALClientDriver = swalClientDriver
	err = p.SWALClientDriver.Init(&p.SoloOSEnv,
		options.SWALAgentPeerID, options.SWALAgentServeAddr,
		options.SWALDBDriver, options.SWALDsn,
		options.SWALDefaultNetBlockCap, options.SWALDefaultMemBlockCap)
	if err != nil {
		return err
	}

	err = p.initSDFS()
	if err != nil {
		return err
	}

	err = p.initSWAL()
	if err != nil {
		return err
	}

	return nil
}

func (p *SoloOS) Serve() error {
	var err error
	err = p.SWALClientDriver.Serve()
	if err != nil {
		return err
	}

	return nil
}

func (p *SoloOS) Close() error {
	var err error
	err = p.SWALClientDriver.Close()
	if err != nil {
		return err
	}

	return nil
}
