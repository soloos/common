package snettypes

import (
	"bytes"
	"soloos/sdbone/offheap"
	"unsafe"

	"soloos/common/util"
)

const (
	PeerStructSize = unsafe.Sizeof(Peer{})
	PeerIDBytesNum = 64
)

type PeerID [PeerIDBytesNum]byte
type PeerUintptr uintptr

func InitTmpPeerID(peerID *PeerID) {
	util.InitUUID64((*[64]byte)(peerID))
}

func StrToPeerID(peerIDStr string) PeerID {
	var ret PeerID
	ret.SetStr(peerIDStr)
	return ret
}

func (p PeerID) Str() string {
	var pos = bytes.IndexByte(p[:], 0)
	if pos == -1 {
		return string(p[:])
	}
	return string(p[:pos])
}

func (p *PeerID) SetStr(peerIDStr string) {
	copy((*p)[:], peerIDStr)
}

func (u PeerUintptr) Ptr() *Peer { return (*Peer)(unsafe.Pointer(u)) }

type PeerJSON struct {
	PeerID          string
	Address         string
	ServiceProtocol int
}

func PeerJSONToPeer(peerJSON PeerJSON) Peer {
	var ret Peer
	ret.ID = StrToPeerID(peerJSON.PeerID)
	ret.SetAddress(peerJSON.Address)
	ret.ServiceProtocol = peerJSON.ServiceProtocol
	return ret
}

func PeerToPeerJSON(peer Peer) PeerJSON {
	var ret PeerJSON
	ret.PeerID = peer.PeerIDStr()
	ret.Address = peer.AddressStr()
	ret.ServiceProtocol = peer.ServiceProtocol
	return ret
}

type Peer struct {
	offheap.LKVTableObjectWithBytes64 `db:"-"`

	addressLen      int
	Address         [128]byte
	ServiceProtocol int
}

func (p *Peer) SetAddressBytes(addr []byte) {
	p.addressLen = len(addr)
	copy(p.Address[:p.addressLen], addr)
}

func (p *Peer) SetAddress(addr string) {
	p.addressLen = len(addr)
	copy(p.Address[:p.addressLen], []byte(addr))
}

func (p *Peer) AddressStr() string {
	return string(p.Address[:p.addressLen])
}

func (p *Peer) SetPeerIDFromStr(peerStr string) {
	copy(p.ID[:], []byte(peerStr))
}

func (p *Peer) PeerID() PeerID { return PeerID(p.ID) }

func (p *Peer) PeerIDStr() string { return PeerID(p.ID).Str() }
