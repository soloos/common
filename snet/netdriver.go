package snet

import (
	"soloos/common/iron"
	"soloos/common/snettypes"
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
		int(snettypes.PeerStructSize), -1, offheap.DefaultKVTableSharedCount*6, nil)
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

func MakeSysPeerID(sysPeerID string) snettypes.PeerID {
	var peerID snettypes.PeerID
	peerID.SetStr("SYS_" + sysPeerID)
	return peerID
}

func (p *NetDriver) InitPeerID(peerID *snettypes.PeerID) {
	// todo: ensure peer id unique
	snettypes.InitTmpPeerID(peerID)
}

func (p *NetDriver) ListPeer(listPeer offheap.LKVTableListObjectWithBytes64) {
	p.peerTable.ListObject(listPeer)
}

func (p *NetDriver) GetPeer(peerID snettypes.PeerID) (snettypes.Peer, error) {
	var uPeer = p.peerTable.TryGetObject(peerID)
	p.peerTable.ReleaseObject(offheap.LKVTableObjectUPtrWithBytes64(uPeer))
	if uPeer == 0 {
		if p.client != nil {
			var peer, err = p.client.GetPeer(peerID)
			if err != nil {
				return snettypes.Peer{}, err
			}
			err = p.doRegisterPeer(peer, true)
			return peer, err
		}
		return snettypes.Peer{}, snettypes.ErrObjectNotExists
	}

	return *snettypes.PeerUintptr(uPeer).Ptr(), nil
}

func (p *NetDriver) doRegisterPeer(peer snettypes.Peer, isSkipRegisterRemote bool) error {
	var (
		uObject        offheap.LKVTableObjectUPtrWithBytes64
		uPeer          snettypes.PeerUintptr
		afterSetNewObj offheap.KVTableAfterSetNewObj
		err            error
	)

	var isNeedUpdateInDB = false

	uObject, afterSetNewObj = p.peerTable.MustGetObject(peer.ID)
	uPeer = snettypes.PeerUintptr(uObject)
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

func (p *NetDriver) RegisterPeer(peer snettypes.Peer) error {
	return p.doRegisterPeer(peer, false)
}

func (p *NetDriver) RegisterPeerInDB(peer snettypes.Peer) error {
	return p.server.RegisterSNetPeerInDB(peer)
}
