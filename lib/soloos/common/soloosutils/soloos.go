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

	err = p.SoloOSEnv.SNetDriver.StartClient(options.SNetDriverServeAddr)
	if err != nil {
		return err
	}

	err = p.initSDFS(sdfsClientDriver)
	if err != nil {
		return err
	}

	err = p.initSWAL(swalClientDriver)
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