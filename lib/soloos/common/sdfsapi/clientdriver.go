package sdfsapi

import (
	"soloos/common/sdfsapitypes"
	"soloos/common/snettypes"
	"soloos/common/soloosbase"
)

type ClientDriver interface {
	Init(soloOSEnv *soloosbase.SoloOSEnv,
		nameNodePeerID snettypes.PeerID,
		dbDriver string, dsn string,
	) error
	InitClient(client Client,
		nameSpaceID sdfsapitypes.NameSpaceID,
		defaultNetBlockCap int,
		defaultMemBlockCap int,
		defaultMemBlocksLimit int32,
	) error
}
