package snet

import (
	"soloos/common/iron"
	"soloos/solodb/offheap"
)

type NetDriver struct {
	offheapDriver *offheap.OffheapDriver
	peerTable     offheap.LKVTableWithBytes64

	server *NetDriverWebServer
	client *NetDriverWebClient
}

func (p *NetDriver) Init(offheapDriver *offheap.OffheapDriver) error {
	var err error
	p.offheapDriver = offheapDriver
	err = p.offheapDriver.InitLKVTableWithBytes64(&p.peerTable, "SNetDriver",
		int(PeerStructSize), -1, offheap.DefaultKVTableSharedCount*6, nil)
	if err != nil {
		return err
	}

	return nil
}

// Serve start NetDriver.NetDriverWebServer
func (p *NetDriver) PrepareServer(webServePrefix string, webServer *iron.Server,
	fetchSNetPeerFromDB FetchSNetPeerFromDB,
	registerSNetPeerInDB RegisterSNetPeerInDB,
) error {
	var err error
	p.server = &NetDriverWebServer{}
	p.server.Init(p,
		webServePrefix, webServer,
		fetchSNetPeerFromDB, registerSNetPeerInDB)
	if err != nil {
		return err
	}

	return nil
}

func (p *NetDriver) ServerServe() error {
	return p.server.Serve()
}

func (p *NetDriver) CloseServer() error {
	return nil
}

// InitClient
func (p *NetDriver) PrepareClient(webServerAddr string) error {
	var err error
	p.client, err = NewNetDriverWebClient(p, webServerAddr)
	if err != nil {
		return err
	}
	return nil
}

func MakeSysPeerID(sysPeerID string) PeerID {
	var peerID PeerID
	peerID.SetStr("SYS_" + sysPeerID)
	return peerID
}

func (p *NetDriver) InitPeerID(peerID *PeerID) {
	// todo: ensure peer id unique
	InitTmpPeerID(peerID)
}

func (p *NetDriver) ListPeer(listPeer offheap.LKVTableListObjectWithBytes64) {
	p.peerTable.ListObject(listPeer)
}

func (p *NetDriver) GetPeer(peerID PeerID) (Peer, error) {
	var uPeer = p.peerTable.TryGetObject(peerID)
	p.peerTable.ReleaseObject(offheap.LKVTableObjectUPtrWithBytes64(uPeer))
	if uPeer == 0 {
		if p.client != nil {
			var peer, err = p.client.GetPeer(peerID)
			if err != nil {
				return Peer{}, err
			}
			err = p.doRegisterPeer(peer, true)
			return peer, err
		}
		return Peer{}, ErrObjectNotExists
	}

	return *PeerUintptr(uPeer).Ptr(), nil
}

func (p *NetDriver) doRegisterPeer(peer Peer, isSkipRegisterRemote bool) error {
	var (
		uObject        offheap.LKVTableObjectUPtrWithBytes64
		uPeer          PeerUintptr
		afterSetNewObj offheap.KVTableAfterSetNewObj
		err            error
	)

	var isNeedUpdateInDB = false

	uObject, afterSetNewObj = p.peerTable.MustGetObject(peer.ID)
	uPeer = PeerUintptr(uObject)
	if afterSetNewObj != nil {
		afterSetNewObj()
		uPeer.Ptr().SetAddress(peer.AddressStr())
		uPeer.Ptr().ServiceProtocol = peer.ServiceProtocol
		isNeedUpdateInDB = true
	} else {
		isNeedUpdateInDB = uPeer.Ptr().Address != peer.Address ||
			uPeer.Ptr().ServiceProtocol != peer.ServiceProtocol
	}
	p.peerTable.ReleaseObject(offheap.LKVTableObjectUPtrWithBytes64(uPeer))

	if p.client != nil && isNeedUpdateInDB && !isSkipRegisterRemote {
		err = p.client.RegisterPeer(peer.ID, peer.AddressStr(), peer.ServiceProtocol)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *NetDriver) RegisterPeer(peer Peer) error {
	return p.doRegisterPeer(peer, false)
}

func (p *NetDriver) RegisterPeerInDB(peer Peer) error {
	return p.server.RegisterSNetPeerInDB(peer)
}
