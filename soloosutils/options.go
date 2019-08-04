package soloosutils

type Options struct {
	SNetDriverServeAddr string
	SoloBoatWebPeerID   string

	SDFSNameNodeSRPCPeerID string
	SDFSDBDriver           string
	SDFSDsn                string

	SWALDefaultNetBlockCap int
	SWALDefaultMemBlockCap int
	SWALBrokerSRPCPeerID   string
	SWALBrokerServeAddr    string
	SWALDBDriver           string
	SWALDsn                string
}
