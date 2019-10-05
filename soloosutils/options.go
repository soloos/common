package soloosutils

type Options struct {
	SNetDriverServeAddr string
	SoloboatWebPeerID   string

	SolofsSolonnSrpcPeerID string
	SolofsDBDriver           string
	SolofsDsn                string

	SolomqDefaultNetBlockCap int
	SolomqDefaultMemBlockCap int
	SolomqSrpcPeerID   string
	SolomqServeAddr    string
	SolomqDBDriver           string
	SolomqDsn                string
}
