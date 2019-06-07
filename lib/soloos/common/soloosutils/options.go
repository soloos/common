package soloosutils

type Options struct {
	SNetDriverServeAddr string

	SDFSNameNodePeerID string
	SDFSDBDriver       string
	SDFSDsn            string

	SWALDefaultNetBlockCap int
	SWALDefaultMemBlockCap int
	SWALAgentPeerID        string
	SWALAgentServeAddr     string
	SWALDBDriver           string
	SWALDsn                string
}
