package sdfsapi

import (
	"soloos/common/soloosbase"
)

type ClientDriver interface {
	Init(soloOSEnv *soloosbase.SoloOSEnv, nameNodeSRPCServerAddr string,
		dbDriver string, dsn string,
	) error
	InitClient(client Client,
		defaultNetBlockCap int,
		defaultMemBlockCap int,
		defaultMemBlocksLimit int32,
	) error
}
