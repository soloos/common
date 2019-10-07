package snet

import (
	"soloos/common/iron"
	"soloos/solodb/offheap"
)

type FetchSNetPeerFromDB func(peerID PeerID) (Peer, error)
type RegisterSNetPeerInDB func(peer Peer) error

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
	var ret []PeerJSON
	p.netDriver.ListPeer(func(uObj offheap.LKVTableObjectUPtrWithBytes64) bool {
		var peer = *PeerUintptr(uObj).Ptr()
		ret = append(ret, PeerToPeerJSON(peer))
		return true
	})
	ir.ApiOutput(ret, CODE_OK, "")
}

func (p *NetDriverWebServer) ctrGetPeer(ir *iron.Request) {
	var (
		peerID PeerID
		peer   Peer
		err    error
	)

	peerID.SetStr(ir.MustFormString("PeerID", ""))
	peer, err = p.netDriver.GetPeer(peerID)
	if err != nil {
		if err == ErrObjectNotExists && p.doFetchSNetPeerFromDB != nil {
			peer, err = p.doFetchSNetPeerFromDB(peerID)
			if err != nil {
				ir.ApiOutput(nil, CODE_502, err.Error())
				return
			}
		}
		return
	}

	ir.ApiOutput(PeerToPeerJSON(peer), CODE_OK, "")
}

func (p *NetDriverWebServer) ctrRegisterPeer(ir *iron.Request) {
	var (
		req RegisterPeerReq
		err error
	)

	err = ir.DecodeBodyJSONData(&req)
	if err != nil {
		ir.ApiOutput(nil, CODE_502, err.Error())
		return
	}

	var peer = Peer{
		LKVTableObjectWithBytes64: offheap.LKVTableObjectWithBytes64{
			ID: StrToPeerID(req.PeerID),
		},
		ServiceProtocol: InitServiceProtocol(req.Protocol),
	}
	peer.SetAddress(req.Addr)

	err = p.RegisterSNetPeerInDB(peer)
	if err != nil {
		ir.ApiOutput(nil, CODE_502, err.Error())
		return
	}

	ir.ApiOutput(nil, CODE_OK, "")
}

func (p *NetDriverWebServer) RegisterSNetPeerInDB(peer Peer) error {
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
