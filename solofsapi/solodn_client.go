package solofsapi

import (
	"soloos/common/snettypes"
	"soloos/common/solofsapitypes"
	"soloos/common/soloosbase"
)

type SolodnClient struct {
	*soloosbase.SoloosEnv
	preadMemBlockWithDisk    solofsapitypes.PReadMemBlockWithDisk
	uploadMemBlockWithDisk   solofsapitypes.UploadMemBlockWithDisk
	uploadMemBlockWithSolomq solofsapitypes.UploadMemBlockWithSolomq
}

func (p *SolodnClient) Init(soloosEnv *soloosbase.SoloosEnv) error {
	p.SoloosEnv = soloosEnv
	return nil
}

func (p *SolodnClient) Dispatch(solodnPeerID snettypes.PeerID,
	path string, ret interface{}, reqArgs ...interface{}) error {
	return p.SNetClientDriver.SimpleCall(solodnPeerID,
		path, ret, reqArgs...)
}
