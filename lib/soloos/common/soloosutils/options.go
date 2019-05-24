package soloosutils

type Options struct {
	SDFSNameNodeServeAddr string
	SDFSDBDriver          string
	SDFSDsn               string

	SWALDefaultNetBlockCap int
	SWALDefaultMemBlockCap int
	SWALAgentPeerID        string
	SWALAgentServeAddr     string
	SWALDBDriver           string
	SWALDsn                string
}
