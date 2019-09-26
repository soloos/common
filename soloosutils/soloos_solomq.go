package soloosutils

import "soloos/common/solomqapi"

func (p *SoloOS) initSOLOMQ(solomqClientDriver solomqapi.ClientDriver) error {
	var err error
	p.SOLOMQClientDriver = solomqClientDriver
	err = p.SOLOMQClientDriver.Init(&p.SoloOSEnv,
		p.options.SoloBoatWebPeerID,
		p.options.SOLOMQBrokerSRPCPeerID, p.options.SOLOMQBrokerServeAddr,
		p.options.SOLOMQDBDriver, p.options.SOLOMQDsn,
		p.options.SOLOMQDefaultNetBlockCap, p.options.SOLOMQDefaultMemBlockCap)
	if err != nil {
		return err
	}

	return nil
}
