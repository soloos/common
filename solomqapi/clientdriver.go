package solomqapi

import (
	"soloos/common/soloosbase"
	"soloos/common/solomqapitypes"
)

type ClientDriver interface {
	Serve() error
	Close() error
	Init(soloOSEnv *soloosbase.SoloOSEnv,
		soloBoatWebPeerID string,
		brokerSRPCPeerIDStr string, brokerSRPCServeAddr string,
		dbDriver string, dsn string,
		defaultNetBlockCap int, defaultMemBlockCap int) error
	InitClient(client Client,
		topicIDStr string, solomqMembers []solomqapitypes.SOLOMQMember) error
}
