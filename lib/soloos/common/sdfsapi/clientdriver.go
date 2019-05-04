package sdfsapi

import (
	soloosbase "soloos/common/soloosapi/base"
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
