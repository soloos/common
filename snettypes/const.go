package snettypes

const (
	SNetVersion           = byte(192)
	SNetMemberCommonLimit = 8
)

const (
	ProtocolUnknown = -1
)

var (
	ProtocolLocalFs  = InitServiceProtocol("localfs")
	ProtocolSolomq   = InitServiceProtocol("solomq")
	ProtocolSolofs   = InitServiceProtocol("solofs")
	ProtocolSrpc     = InitServiceProtocol("srpc")
	ProtocolSoloboat = InitServiceProtocol("soloboat")
	ProtocolWeb      = InitServiceProtocol("web")
)
