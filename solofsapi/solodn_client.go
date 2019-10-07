package solofsapi

import (
	"soloos/common/snet"
	"soloos/common/solofstypes"
	"soloos/common/soloosbase"
)

type SolodnClient struct {
	*soloosbase.SoloosEnv
	preadMemBlockWithDisk    solofstypes.PReadMemBlockWithDisk
	uploadMemBlockWithDisk   solofstypes.UploadMemBlockWithDisk
	uploadMemBlockWithSolomq solofstypes.UploadMemBlockWithSolomq
}

func (p *SolodnClient) Init(soloosEnv *soloosbase.SoloosEnv) error {
	p.SoloosEnv = soloosEnv
	return nil
}

func (p *SolodnClient) Dispatch(solodnPeerID snet.PeerID,
	path string, ret interface{}, reqArgs ...interface{}) error {
	return p.SimpleCall(solodnPeerID,
		path, ret, reqArgs...)
}
