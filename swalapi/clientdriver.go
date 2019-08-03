package swalapi

import (
	"soloos/common/soloosbase"
	"soloos/common/swalapitypes"
)

type ClientDriver interface {
	Serve() error
	Close() error
	Init(soloOSEnv *soloosbase.SoloOSEnv,
		brokerPeerIDStr string, brokerServeAddr string,
		dbDriver string, dsn string,
		defaultNetBlockCap int, defaultMemBlockCap int) error
	InitClient(client Client,
		topicIDStr string, swalMembers []swalapitypes.SWALMember) error
}
