package soloosutils

import (
	"soloos/common/iron"
	"soloos/common/log"
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

var _ = iron.IServer(&SoloOS{})

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
	err = p.SoloOSEnv.InitWithSNet(options.SNetDriverServeAddr)
	if err != nil {
		log.Warn("SoloOSEnv Init error", err)
		return err
	}

	err = p.initSDFS(sdfsClientDriver)
	if err != nil {
		log.Warn("SoloOS initSDFS error", err)
		return err
	}

	err = p.initSWAL(swalClientDriver)
	if err != nil {
		log.Warn("SoloOS initSWAL error", err)
		return err
	}

	return nil
}

func (p *SoloOS) ServerName() string {
	return "SoloOS.Instance"
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
