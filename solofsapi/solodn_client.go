package solofsapi

import (
	"soloos/common/solofsapitypes"
	"soloos/common/soloosbase"
)

type SolodnClient struct {
	*soloosbase.SoloOSEnv
	preadMemBlockWithDisk  solofsapitypes.PReadMemBlockWithDisk
	uploadMemBlockWithDisk solofsapitypes.UploadMemBlockWithDisk
	uploadMemBlockWithSOLOMQ solofsapitypes.UploadMemBlockWithSOLOMQ
}

func (p *SolodnClient) Init(soloOSEnv *soloosbase.SoloOSEnv) error {
	p.SoloOSEnv = soloOSEnv
	return nil
}

func (p *SolodnClient) SetPReadMemBlockWithDisk(preadMemBlockWithDisk solofsapitypes.PReadMemBlockWithDisk) {
	p.preadMemBlockWithDisk = preadMemBlockWithDisk
}

func (p *SolodnClient) SetUploadMemBlockWithDisk(uploadMemBlockWithDisk solofsapitypes.UploadMemBlockWithDisk) {
	p.uploadMemBlockWithDisk = uploadMemBlockWithDisk
}

func (p *SolodnClient) SetUploadMemBlockWithSOLOMQ(uploadMemBlockWithSOLOMQ solofsapitypes.UploadMemBlockWithSOLOMQ) {
	p.uploadMemBlockWithSOLOMQ = uploadMemBlockWithSOLOMQ
}
