package snettypes

const (
	SNetVersion = byte(192)
)

const (
	ProtocolUnknown = -1
)

var (
	ProtocolLocalFS     = InitServiceProtocol("localfs")
	ProtocolSOLOMQ     = InitServiceProtocol("solomq")
	ProtocolSOLOFS     = InitServiceProtocol("solofs")
	ProtocolSRPC     = InitServiceProtocol("srpc")
	ProtocolSoloBoat = InitServiceProtocol("soloboat")
	ProtocolWeb      = InitServiceProtocol("web")
)
