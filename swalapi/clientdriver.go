package swalapi

import (
	"soloos/common/soloosbase"
	"soloos/common/swalapitypes"
)

type ClientDriver interface {
	Serve() error
	Close() error
	Init(soloOSEnv *soloosbase.SoloOSEnv,
		brokerSRPCPeerIDStr string, brokerSRPCServeAddr string,
		dbDriver string, dsn string,
		defaultNetBlockCap int, defaultMemBlockCap int) error
	InitClient(client Client,
		topicIDStr string, swalMembers []swalapitypes.SWALMember) error
}
