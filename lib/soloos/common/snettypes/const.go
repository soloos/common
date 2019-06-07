package snettypes

const (
	SNetVersion = byte(192)
)

const (
	ProtocolUnknown = -1
	ProtocolDisk    = iota
	ProtocolSWAL
	ProtocolSRPC
)
