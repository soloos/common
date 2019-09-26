package soloosutils

type Options struct {
	SNetDriverServeAddr string
	SoloBoatWebPeerID   string

	SOLOFSSolonnSRPCPeerID string
	SOLOFSDBDriver           string
	SOLOFSDsn                string

	SOLOMQDefaultNetBlockCap int
	SOLOMQDefaultMemBlockCap int
	SOLOMQBrokerSRPCPeerID   string
	SOLOMQBrokerServeAddr    string
	SOLOMQDBDriver           string
	SOLOMQDsn                string
}
