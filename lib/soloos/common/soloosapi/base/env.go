package base

import (
	"soloos/common/snet"
	"soloos/sdbone/offheap"
)

type SoloOSEnv struct {
	OffheapDriver    offheap.OffheapDriver
	SNetDriver       snet.NetDriver
	SNetClientDriver snet.ClientDriver
}

func (p *SoloOSEnv) Init() error {
	var err error

	err = p.OffheapDriver.Init()
	if err != nil {
		return err
	}

	err = p.SNetDriver.Init(&p.OffheapDriver, "SoloOSEnv")
	if err != nil {
		return err
	}

	err = p.SNetClientDriver.Init(&p.OffheapDriver)
	if err != nil {
		return err
	}

	return nil
}