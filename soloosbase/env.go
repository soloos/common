package soloosbase

import (
	"soloos/common/snet"
	"soloos/sdbone/offheap"
)

type SoloOSEnv struct {
	OffheapDriver    offheap.OffheapDriver
	SNetDriver       snet.NetDriver
	SNetClientDriver snet.SRPCClientDriver
}

func (p *SoloOSEnv) InitWithSNet(snetWebServerAddr string) error {
	var err error

	err = p.OffheapDriver.Init()
	if err != nil {
		return err
	}

	err = p.SNetDriver.Init(&p.OffheapDriver)
	if err != nil {
		return err
	}

	err = p.SNetClientDriver.Init(&p.OffheapDriver, &p.SNetDriver)
	if err != nil {
		return err
	}

	if snetWebServerAddr != "" {
		err = p.SNetDriver.PrepareClient(snetWebServerAddr)
		if err != nil {
			return err
		}
	}

	return nil
}
