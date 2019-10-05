package solofsapi

import (
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

func (p *SolodnClient) SetPReadMemBlockWithDisk(preadMemBlockWithDisk solofsapitypes.PReadMemBlockWithDisk) {
	p.preadMemBlockWithDisk = preadMemBlockWithDisk
}

func (p *SolodnClient) SetUploadMemBlockWithDisk(uploadMemBlockWithDisk solofsapitypes.UploadMemBlockWithDisk) {
	p.uploadMemBlockWithDisk = uploadMemBlockWithDisk
}

func (p *SolodnClient) SetUploadMemBlockWithSolomq(uploadMemBlockWithSolomq solofsapitypes.UploadMemBlockWithSolomq) {
	p.uploadMemBlockWithSolomq = uploadMemBlockWithSolomq
}
