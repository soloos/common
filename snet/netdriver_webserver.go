package snet

import (
	"soloos/common/iron"
	"soloos/common/snettypes"
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

	p.server.Router("/Peer/List", p.ctrPeerList)
	p.server.Router("/Peer/Get", p.ctrGetPeer)
	p.server.Router("/Peer/Register", p.ctrRegisterPeer)

	return nil
}

func (p *NetDriverWebServer) Serve() error {
	return p.server.Serve()
}

func (p *NetDriverWebServer) ctrPeerList(ir *iron.Request) {
	var ret []snettypes.PeerJSON
	p.netDriver.ListPeer(func(uObj offheap.LKVTableObjectUPtrWithBytes64) bool {
		var peer = *snettypes.PeerUintptr(uObj).Ptr()
		ret = append(ret, snettypes.PeerToPeerJSON(peer))
		return true
	})
	ir.ApiOutput(ret, snettypes.CODE_OK, "")
}

func (p *NetDriverWebServer) ctrGetPeer(ir *iron.Request) {
	var (
		peerID snettypes.PeerID
		req    GetPeerReqJSON
		peer   snettypes.Peer
		err    error
	)
	err = ir.DecodeBodyJSONData(&req)
	if err != nil {
		ir.ApiOutput(nil, snettypes.CODE_502, err.Error())
		return
	}

	peerID.SetStr(req.PeerID)
	peer, err = p.netDriver.GetPeer(peerID)
	if err != nil {
		if err == snettypes.ErrObjectNotExists && p.FetchSNetPeerFromDB != nil {
			peer, err = p.FetchSNetPeerFromDB(peerID)
			if err != nil {
				ir.ApiOutput(nil, snettypes.CODE_502, err.Error())
				return
			}
		}
		ir.ApiOutput(nil, snettypes.CODE_502, err.Error())
		return
	}

	ir.ApiOutput(snettypes.PeerToPeerJSON(peer), snettypes.CODE_OK, "")
}

func (p *NetDriverWebServer) ctrRegisterPeer(ir *iron.Request) {
	var (
		req RegisterPeerReqJSON
		err error
	)

	err = ir.DecodeBodyJSONData(&req)
	if err != nil {
		ir.ApiOutput(nil, snettypes.CODE_502, err.Error())
		return
	}

	var peer = snettypes.Peer{
		LKVTableObjectWithBytes64: offheap.LKVTableObjectWithBytes64{
			ID: snettypes.StrToPeerID(req.PeerID),
		},
		ServiceProtocol: snettypes.InitServiceProtocol(req.Protocol),
	}
	peer.SetAddress(req.Addr)

	err = p.netDriver.RegisterPeer(peer)
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

	ir.ApiOutput(nil, snettypes.CODE_OK, "")
}
