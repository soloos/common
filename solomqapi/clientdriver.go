package solomqapi

import (
	"soloos/common/soloosbase"
	"soloos/common/solomqapitypes"
)

type ClientDriver interface {
	Serve() error
	Close() error
	Init(soloosEnv *soloosbase.SoloosEnv,
		soloBoatWebPeerID string,
		solomqSRPCPeerIDStr string, solomqSRPCServeAddr string,
		dbDriver string, dsn string,
		defaultNetBlockCap int, defaultMemBlockCap int) error
	InitClient(client Client,
		topicIDStr string, solomqMembers []solomqapitypes.SolomqMember) error
}
