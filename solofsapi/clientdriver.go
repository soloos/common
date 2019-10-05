package solofsapi

import (
	"soloos/common/solofsapitypes"
	"soloos/common/snettypes"
	"soloos/common/soloosbase"
)

type ClientDriver interface {
	Init(soloosEnv *soloosbase.SoloosEnv,
		solonnSrpcPeerID snettypes.PeerID,
		dbDriver string, dsn string,
	) error
	InitClient(client Client,
		nameSpaceID solofsapitypes.NameSpaceID,
		defaultNetBlockCap int,
		defaultMemBlockCap int,
		defaultMemBlocksLimit int32,
	) error
}
