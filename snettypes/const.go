package snettypes

const (
	SNetVersion = byte(192)
)

const (
	ProtocolUnknown = -1
)

var (
	ProtocolLocalFS     = InitServiceProtocol("localfs")
	ProtocolSWAL     = InitServiceProtocol("swal")
	ProtocolSDFS     = InitServiceProtocol("sdfs")
	ProtocolSRPC     = InitServiceProtocol("srpc")
	ProtocolSoloBoat = InitServiceProtocol("soloboat")
	ProtocolWeb      = InitServiceProtocol("web")
)
