package base

type Options struct {
	SDFSNameNodeServeAddr string
	SDFSDBDriver          string
	SDFSDsn               string

	SWALAgentPeerID    string
	SWALAgentServeAddr string
	SWALDBDriver       string
	SWALDsn            string
}
