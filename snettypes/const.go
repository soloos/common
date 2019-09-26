package snettypes

const (
	SNetVersion = byte(192)
)

const (
	ProtocolUnknown = -1
)

var (
	ProtocolLocalFS     = InitServiceProtocol("localfs")
	ProtocolSolomq     = InitServiceProtocol("solomq")
	ProtocolSolofs     = InitServiceProtocol("solofs")
	ProtocolSRPC     = InitServiceProtocol("srpc")
	ProtocolSoloboat = InitServiceProtocol("soloboat")
	ProtocolWeb      = InitServiceProtocol("web")
)
