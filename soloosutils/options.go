package soloosutils

type Options struct {
	SNetDriverServeAddr string
	SoloboatWebPeerID   string

	SolofsSolonnSRPCPeerID string
	SolofsDBDriver           string
	SolofsDsn                string

	SolomqDefaultNetBlockCap int
	SolomqDefaultMemBlockCap int
	SolomqSRPCPeerID   string
	SolomqServeAddr    string
	SolomqDBDriver           string
	SolomqDsn                string
}
