package soloosutils

import (
	"soloos/common/iron"
	"soloos/common/log"
	"soloos/common/solomqapi"
	"soloos/common/solofsapi"
	"soloos/common/soloosbase"
)

var (
	SoloOSInstance        SoloOS
	isDefaultSoloOSInited bool = false
)

type SoloOS struct {
	options Options
	soloosbase.SoloOSEnv

	SOLOFSClientDriver solofsapi.ClientDriver
	SOLOMQClientDriver solomqapi.ClientDriver
}

var _ = iron.IServer(&SoloOS{})

func InitSoloOSInstance(options Options,
	solofsClientDriver solofsapi.ClientDriver,
	solomqClientDriver solomqapi.ClientDriver,
) error {
	if isDefaultSoloOSInited {
		return nil
	}
	isDefaultSoloOSInited = true
	return SoloOSInstance.Init(options, solofsClientDriver, solomqClientDriver)
}

func (p *SoloOS) Init(options Options,
	solofsClientDriver solofsapi.ClientDriver,
	solomqClientDriver solomqapi.ClientDriver,
) error {
	var err error

	p.options = options
	err = p.SoloOSEnv.InitWithSNet(options.SNetDriverServeAddr)
	if err != nil {
		log.Warn("SoloOSEnv Init error", err)
		return err
	}

	err = p.initSOLOFS(solofsClientDriver)
	if err != nil {
		log.Warn("SoloOS initSOLOFS error", err)
		return err
	}

	err = p.initSOLOMQ(solomqClientDriver)
	if err != nil {
		log.Warn("SoloOS initSOLOMQ error", err)
		return err
	}

	return nil
}

func (p *SoloOS) ServerName() string {
	return "SoloOS.Instance"
}

func (p *SoloOS) Serve() error {
	var err error
	err = p.SOLOMQClientDriver.Serve()
	if err != nil {
		return err
	}

	return nil
}

func (p *SoloOS) Close() error {
	var err error
	err = p.SOLOMQClientDriver.Close()
	if err != nil {
		return err
	}

	return nil
}
