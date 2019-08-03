package soloosutils

type Options struct {
	SNetDriverServeAddr string

	SDFSNameNodePeerID string
	SDFSDBDriver       string
	SDFSDsn            string

	SWALDefaultNetBlockCap int
	SWALDefaultMemBlockCap int
	SWALBrokerPeerID        string
	SWALBrokerServeAddr     string
	SWALDBDriver           string
	SWALDsn                string
}
