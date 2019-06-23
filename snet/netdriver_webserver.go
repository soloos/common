package snet

import (
	"soloos/common/snettypes"
	"soloos/common/iron"
	"soloos/sdbone/offheap"
)

type FetchSNetPeerFromDB func(peerID snettypes.PeerID) (snettypes.Peer, error)
type RegisterSNetPeerInDB func(peer snettypes.Peer) error

type NetDriverWebServer struct {
	netDriver   *NetDriver
	webServeStr string
	server      iron.Server

	FetchSNetPeerFromDB  FetchSNetPeerFromDB
	RegisterSNetPeerInDB RegisterSNetPeerInDB
}

func NewNetDriverWebServer(netDriver *NetDriver,
	webServeAddr string,
	fetchSNetPeerFromDB FetchSNetPeerFromDB,
	registerSNetPeerInDB RegisterSNetPeerInDB,
	webOptions iron.Options) (*NetDriverWebServer, error) {
	var (
		ret *NetDriverWebServer = new(NetDriverWebServer)
		err error
	)

	ret.Init(netDriver,
		webServeAddr,
		fetchSNetPeerFromDB,
		registerSNetPeerInDB,
		webOptions)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (p *NetDriverWebServer) Init(netDriver *NetDriver,
	webServeAddr string,
	fetchSNetPeerFromDB FetchSNetPeerFromDB,
	registerSNetPeerInDB RegisterSNetPeerInDB,
	webOptions iron.Options) error {
	var err error

	p.netDriver = netDriver
	p.webServeStr = webServeAddr
	p.FetchSNetPeerFromDB = fetchSNetPeerFromDB
	p.RegisterSNetPeerInDB = registerSNetPeerInDB

	err = p.server.Init(webOptions)
	if err != nil {
		return err
	}

	p.server.Router("/Peer/Get", p.ctrGetPeer)
	p.server.Router("/Peer/Register", p.ctrRegisterPeer)

	return nil
}

func (p *NetDriverWebServer) Serve() error {
	return p.server.Serve()
}

func (p *NetDriverWebServer) ctrGetPeer(ir *iron.Request) {
	var peerID snettypes.PeerID
	peerID.SetStr(ir.MustFormString("PeerID", ""))
	var ret, err = p.netDriver.GetPeer(peerID)
	if err != nil {
		if err == snettypes.ErrObjectNotExists && p.FetchSNetPeerFromDB != nil {
			ret, err = p.FetchSNetPeerFromDB(peerID)
			if err != nil {
				ir.ApiOutput(nil, snettypes.CODE_502, err.Error())
				return
			}
		}
		ir.ApiOutput(nil, snettypes.CODE_502, err.Error())
		return
	}

	ir.ApiOutput(snettypes.PeerToPeerJSON(ret), snettypes.CODE_OK, "")
}

func (p *NetDriverWebServer) ctrRegisterPeer(ir *iron.Request) {
	var peer = snettypes.Peer{
		LKVTableObjectWithBytes64: offheap.LKVTableObjectWithBytes64{
			ID: snettypes.StrToPeerID(ir.MustFormString("PeerID", "")),
		},
		ServiceProtocol: ir.MustFormInt("Protocol", snettypes.ProtocolUnknown),
	}
	peer.SetAddress(ir.MustFormString("Addr", ""))

	var err = p.netDriver.RegisterPeer(peer)
	if err != nil {
		ir.ApiOutput(nil, snettypes.CODE_502, err.Error())
		return
	}

	if p.RegisterSNetPeerInDB != nil {
		err = p.RegisterSNetPeerInDB(peer)
		if err != nil {
			ir.ApiOutput(nil, snettypes.CODE_502, err.Error())
			return
		}
	}

	ir.ApiOutput(snettypes.PeerToPeerJSON(peer), snettypes.CODE_OK, "")
}
