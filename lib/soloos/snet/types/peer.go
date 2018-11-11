package types

import (
	"unsafe"
)

const (
	PeerStructSize = unsafe.Sizeof(Peer{})
)

type PeerID = [64]byte

type PeerUintptr uintptr

func (u PeerUintptr) Ptr() *Peer { return (*Peer)(unsafe.Pointer(u)) }

type Peer struct {
	ID                 PeerID
	Address            [128]byte
	ServiceProtocol    [32]byte
	addressLen         int
	serviceProtocolLen int
}

func (p *Peer) SetAddress(addr string) {
	p.addressLen = len(addr)
	copy(p.Address[:p.addressLen], []byte(addr))
}

func (p *Peer) AddressStr() string {
	return string(p.Address[:p.addressLen])
}

func (p *Peer) SetServiceProtocol(protocol string) {
	p.serviceProtocolLen = len(protocol)
	copy(p.ServiceProtocol[:p.serviceProtocolLen], []byte(protocol))
}

func (p *Peer) ServiceProtocolStr() string {
	return string(p.ServiceProtocol[:p.serviceProtocolLen])
}

func InitPeerID(peerID *PeerID) {
}
