package snet

import (
	"soloos/common/iron"
	"soloos/common/snettypes"
	"soloos/sdbone/offheap"
)

type FetchSNetPeerFromDB func(peerID snettypes.PeerID) (snettypes.Peer, error)
type RegisterSNetPeerInDB func(peer snettypes.Peer) error

type NetDriverWebServer struct {
	netDriver      *NetDriver
	webServePrefix string
	server         *iron.Server

	doFetchSNetPeerFromDB  FetchSNetPeerFromDB
	doRegisterSNetPeerInDB RegisterSNetPeerInDB
}

func (p *NetDriverWebServer) Init(netDriver *NetDriver,
	webServePrefix string, webServer *iron.Server,
	fetchSNetPeerFromDB FetchSNetPeerFromDB,
	registerSNetPeerInDB RegisterSNetPeerInDB) error {

	p.netDriver = netDriver
	p.webServePrefix = webServePrefix
	p.server = webServer
	p.doFetchSNetPeerFromDB = fetchSNetPeerFromDB
	p.doRegisterSNetPeerInDB = registerSNetPeerInDB

	p.server.Router(p.webServePrefix+"/Peer/List", p.ctrPeerList)
	p.server.Router(p.webServePrefix+"/Peer/Get", p.ctrGetPeer)
	p.server.Router(p.webServePrefix+"/Peer/Register", p.ctrRegisterPeer)

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
		peer   snettypes.Peer
		err    error
	)

	peerID.SetStr(ir.MustFormString("PeerID", ""))
	peer, err = p.netDriver.GetPeer(peerID)
	if err != nil {
		if err == snettypes.ErrObjectNotExists && p.doFetchSNetPeerFromDB != nil {
			peer, err = p.doFetchSNetPeerFromDB(peerID)
			if err != nil {
				ir.ApiOutput(nil, snettypes.CODE_502, err.Error())
				return
			}
		}
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

	err = p.RegisterSNetPeerInDB(peer)
	if err != nil {
		ir.ApiOutput(nil, snettypes.CODE_502, err.Error())
		return
	}

	ir.ApiOutput(nil, snettypes.CODE_OK, "")
}

func (p *NetDriverWebServer) RegisterSNetPeerInDB(peer snettypes.Peer) error {
	var err error
	err = p.netDriver.RegisterPeer(peer)
	if err != nil {
		return err
	}

	if p.doRegisterSNetPeerInDB != nil {
		err = p.doRegisterSNetPeerInDB(peer)
		if err != nil {
			return err
		}
	}

	return nil
}
