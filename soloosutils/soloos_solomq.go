package soloosutils

import "soloos/common/solomqapi"

func (p *Soloos) initSolomq(solomqClientDriver solomqapi.ClientDriver) error {
	var err error
	p.SolomqClientDriver = solomqClientDriver
	err = p.SolomqClientDriver.Init(&p.SoloosEnv,
		p.options.SoloboatWebPeerID,
		p.options.SolomqSRPCPeerID, p.options.SolomqServeAddr,
		p.options.SolomqDBDriver, p.options.SolomqDsn,
		p.options.SolomqDefaultNetBlockCap, p.options.SolomqDefaultMemBlockCap)
	if err != nil {
		return err
	}

	return nil
}
