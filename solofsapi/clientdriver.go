package solofsapi

import (
	"soloos/common/snettypes"
	"soloos/common/solofsapitypes"
	"soloos/common/soloosbase"
)

type ClientDriver interface {
	Init(soloosEnv *soloosbase.SoloosEnv,
		solonnSrpcPeerID snettypes.PeerID,
		dbDriver string, dsn string,
	) error
	InitClient(client Client,
		nsID solofsapitypes.NameSpaceID,
		defaultNetBlockCap int,
		defaultMemBlockCap int,
		defaultMemBlocksLimit int32,
	) error
}
