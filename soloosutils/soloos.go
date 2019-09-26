package soloosutils

import (
	"soloos/common/iron"
	"soloos/common/log"
	"soloos/common/solomqapi"
	"soloos/common/solofsapi"
	"soloos/common/soloosbase"
)

var (
	SoloosInstance        Soloos
	isDefaultSoloosInited bool = false
)

type Soloos struct {
	options Options
	soloosbase.SoloosEnv

	SolofsClientDriver solofsapi.ClientDriver
	SolomqClientDriver solomqapi.ClientDriver
}

var _ = iron.IServer(&Soloos{})

func InitSoloosInstance(options Options,
	solofsClientDriver solofsapi.ClientDriver,
	solomqClientDriver solomqapi.ClientDriver,
) error {
	if isDefaultSoloosInited {
		return nil
	}
	isDefaultSoloosInited = true
	return SoloosInstance.Init(options, solofsClientDriver, solomqClientDriver)
}

func (p *Soloos) Init(options Options,
	solofsClientDriver solofsapi.ClientDriver,
	solomqClientDriver solomqapi.ClientDriver,
) error {
	var err error

	p.options = options
	err = p.SoloosEnv.InitWithSNet(options.SNetDriverServeAddr)
	if err != nil {
		log.Warn("SoloosEnv Init error", err)
		return err
	}

	err = p.initSolofs(solofsClientDriver)
	if err != nil {
		log.Warn("Soloos initSolofs error", err)
		return err
	}

	err = p.initSolomq(solomqClientDriver)
	if err != nil {
		log.Warn("Soloos initSolomq error", err)
		return err
	}

	return nil
}

func (p *Soloos) ServerName() string {
	return "Soloos.Instance"
}

func (p *Soloos) Serve() error {
	var err error
	err = p.SolomqClientDriver.Serve()
	if err != nil {
		return err
	}

	return nil
}

func (p *Soloos) Close() error {
	var err error
	err = p.SolomqClientDriver.Close()
	if err != nil {
		return err
	}

	return nil
}
