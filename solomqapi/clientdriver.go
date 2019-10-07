package solomqapi

import (
	"soloos/common/solomqtypes"
	"soloos/common/soloosbase"
)

type ClientDriver interface {
	Serve() error
	Close() error
	Init(soloosEnv *soloosbase.SoloosEnv,
		soloBoatWebPeerID string,
		solomqSrpcPeerIDStr string, solomqSrpcServeAddr string,
		dbDriver string, dsn string,
		defaultNetBlockCap int, defaultMemBlockCap int) error
	InitClient(client Client,
		topicIDStr string, solomqMembers []solomqtypes.SolomqMember) error
}
