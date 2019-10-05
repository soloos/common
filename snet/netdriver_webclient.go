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
		urlPath = p.webServerAddr + "/Peer/Get?PeerID=" + peerID.Str()
		resp    GetPeerResp
		err     error
	)

	err = iron.HttpGetJSON(urlPath,
		&resp)
	if err != nil {
		return ret, err
	}

	if resp.Error != "" {
		return ret, xerrors.New(resp.Error)
	}

	ret = snettypes.PeerJSONToPeer(resp.Data)
	return ret, nil
}

func (p *NetDriverWebClient) RegisterPeer(peerID snettypes.PeerID, addr string, protocol snettypes.ServiceProtocol) error {
	var (
		urlPath = p.webServerAddr + "/Peer/Register"
		resp    RegisterPeerResp
		err     error
	)

	switch protocol {
	case snettypes.ProtocolLocalFs:
		return nil
	default:
	}

	err = iron.PostJSON(urlPath,
		RegisterPeerReq{
			PeerID:   peerID.Str(),
			Addr:     addr,
			Protocol: protocol.Str(),
		},
		&resp)
	if err != nil {
		return err
	}

	if resp.Error != "" {
		return xerrors.New(resp.Error)
	}

	return nil
}
