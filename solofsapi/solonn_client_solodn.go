package solofsapi

import (
	"soloos/common/snettypes"
	"soloos/common/solofsprotocol"
)

func (p *SolonnClient) SolodnRegister(peerID snettypes.PeerID,
	serveAddr string,
	protocolType snettypes.ServiceProtocol) error {
	var req = solofsprotocol.SNetPeer{
		PeerID:   peerID.Str(),
		Address:  serveAddr,
		Protocol: protocolType.Str(),
	}

	return p.SNetClientDriver.SimpleCall(p.solonnPeerID,
		"/Solodn/Register", req, nil)
}
