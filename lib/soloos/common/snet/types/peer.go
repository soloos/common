package types

import (
	"soloos/sdbone/offheap"
	"unsafe"
)

const (
	PeerStructSize = unsafe.Sizeof(Peer{})
	PeerIDBytesNum = 64
)

type PeerID = [PeerIDBytesNum]byte

type PeerUintptr uintptr

func (u PeerUintptr) Ptr() *Peer { return (*Peer)(unsafe.Pointer(u)) }

type Peer struct {
	offheap.LKVTableObjectWithBytes64 `db:"-"`

	addressLen      int
	Address         [128]byte
	ServiceProtocol int
}

func (p *Peer) SetAddress(addr string) {
	p.addressLen = len(addr)
	copy(p.Address[:p.addressLen], []byte(addr))
}

func (p *Peer) AddressStr() string {
	return string(p.Address[:p.addressLen])
}

func (p *Peer) PeerIDStr() string { return string(p.ID[:]) }
