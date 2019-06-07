package sdfsapi

import (
	"soloos/common/snettypes"
	"soloos/common/soloosbase"
)

type ClientDriver interface {
	Init(soloOSEnv *soloosbase.SoloOSEnv,
		nameNodePeerID snettypes.PeerID,
		dbDriver string, dsn string,
	) error
	InitClient(client Client,
		defaultNetBlockCap int,
		defaultMemBlockCap int,
		defaultMemBlocksLimit int32,
	) error
}
