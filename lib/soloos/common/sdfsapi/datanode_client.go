package sdfsapi

import (
	"soloos/common/sdfsapitypes"
	"soloos/common/soloosbase"
)

type DataNodeClient struct {
	*soloosbase.SoloOSEnv
	preadMemBlockWithDisk  sdfsapitypes.PReadMemBlockWithDisk
	uploadMemBlockWithDisk sdfsapitypes.UploadMemBlockWithDisk
	uploadMemBlockWithSWAL sdfsapitypes.UploadMemBlockWithSWAL
}

func (p *DataNodeClient) Init(soloOSEnv *soloosbase.SoloOSEnv) error {
	p.SoloOSEnv = soloOSEnv
	return nil
}

func (p *DataNodeClient) SetPReadMemBlockWithDisk(preadMemBlockWithDisk sdfsapitypes.PReadMemBlockWithDisk) {
	p.preadMemBlockWithDisk = preadMemBlockWithDisk
}

func (p *DataNodeClient) SetUploadMemBlockWithDisk(uploadMemBlockWithDisk sdfsapitypes.UploadMemBlockWithDisk) {
	p.uploadMemBlockWithDisk = uploadMemBlockWithDisk
}

func (p *DataNodeClient) SetUploadMemBlockWithSWAL(uploadMemBlockWithSWAL sdfsapitypes.UploadMemBlockWithSWAL) {
	p.uploadMemBlockWithSWAL = uploadMemBlockWithSWAL
}
