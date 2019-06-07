package soloosutils

import (
	"soloos/common/sdfsapi"
	"soloos/common/snettypes"
)

func (p *SoloOS) initSDFS(sdfsClientDriver sdfsapi.ClientDriver) error {
	var err error
	var nameNodePeerID snettypes.PeerID
	nameNodePeerID.SetStr(p.options.SDFSNameNodePeerID)
	p.SDFSClientDriver = sdfsClientDriver
	err = p.SDFSClientDriver.Init(&p.SoloOSEnv,
		nameNodePeerID,
		p.options.SDFSDBDriver, p.options.SDFSDsn)
	if err != nil {
		return err
	}

	return nil
}
