package snet

import (
	"soloos/common/iron"
	"soloos/common/snettypes"
	"soloos/common/xerrors"
)

type NetDriverWebClient struct {
	netDriver     *NetDriver
	webServerAddr string
}

func NewNetDriverWebClient(netDriver *NetDriver, webServerAddr string) (*NetDriverWebClient, error) {
	var (
		ret *NetDriverWebClient = new(NetDriverWebClient)
		err error
	)

	ret.Init(netDriver, webServerAddr)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (p *NetDriverWebClient) Init(netDriver *NetDriver, webServerAddr string) error {
	p.webServerAddr = webServerAddr
	return nil
}

func (p *NetDriverWebClient) GetPeer(peerID snettypes.PeerID) (snettypes.Peer, error) {
	var (
		ret     snettypes.Peer
		urlPath = p.webServerAddr + "/Peer/Get"
		resp    GetPeerRespJSON
		err     error
	)

	err = iron.PostJSON(urlPath,
		GetPeerReqJSON{PeerID: peerID.Str()},
		&resp)
	if err != nil {
		return ret, err
	}

	if resp.Errno != snettypes.CODE_OK {
		return ret, xerrors.New(resp.ErrMsg)
	}

	ret = snettypes.PeerJSONToPeer(resp.Data)
	return ret, nil
}

// MustGetPee return uPeer and peer is inited before
func (p *NetDriverWebClient) RegisterPeer(peerID snettypes.PeerID, addr string, protocol snettypes.ServiceProtocol) error {
	var (
		urlPath = p.webServerAddr + "/Peer/Register"
		resp    RegisterPeerRespJSON
		err     error
	)

	switch protocol {
	case snettypes.ProtocolDisk:
		return nil
	default:
	}

	err = iron.PostJSON(urlPath,
		RegisterPeerReqJSON{
			PeerID:   peerID.Str(),
			Addr:     addr,
			Protocol: protocol.Str(),
		},
		&resp)
	if err != nil {
		return err
	}

	if resp.Errno != snettypes.CODE_OK {
		return xerrors.New(resp.ErrMsg)
	}

	return nil
}
