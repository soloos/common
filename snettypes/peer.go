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
type ServiceProtocol [8]byte

func InitServiceProtocol(str string) ServiceProtocol {
	var ret ServiceProtocol
	copy(ret[:], []byte(str))
	return ret
}

func (p *ServiceProtocol) SetProtocolStr(str string) {
	for i, _ := range p {
		p[i] = 0
	}
	copy(p[:], []byte(str))
}

func (p *ServiceProtocol) SetProtocolBytes(bytes []byte) {
	copy(p[:], bytes)
}

func (p *ServiceProtocol) Str() string {
	n := -1
	for i, b := range p {
		if b == 0 {
			break
		}
		n = i
	}
	return string(p[:n+1])
}

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
	ServiceProtocol string
}

func PeerJSONToPeer(peerJSON PeerJSON) Peer {
	var ret Peer
	ret.ID = StrToPeerID(peerJSON.PeerID)
	ret.SetAddress(peerJSON.Address)
	ret.ServiceProtocol.SetProtocolStr(peerJSON.ServiceProtocol)
	return ret
}

func PeerToPeerJSON(peer Peer) PeerJSON {
	var ret PeerJSON
	ret.PeerID = peer.PeerIDStr()
	ret.Address = peer.AddressStr()
	ret.ServiceProtocol = peer.ServiceProtocol.Str()
	return ret
}

type Peer struct {
	offheap.LKVTableObjectWithBytes64 `db:"-"`

	addressLen      int
	Address         [128]byte
	ServiceProtocol ServiceProtocol
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
