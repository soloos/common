package snettypes

const (
	SNetVersion = byte(192)
)

const (
	ProtocolUnknown = -1
)

var (
	ProtocolDisk    = InitServiceProtocol("disk")
	ProtocolSWAL    = InitServiceProtocol("swal")
	ProtocolSDFS    = InitServiceProtocol("sdfs")
	ProtocolSRPC    = InitServiceProtocol("srpc")
	ProtocolSilicon = InitServiceProtocol("silicon")
)
