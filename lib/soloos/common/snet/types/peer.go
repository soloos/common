package types

import (
	"sync"
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
	PeerID           PeerID
	addressLen       int
	Address          [128]byte
	ServiceProtocol  int
	MetaDataMutex    sync.Mutex `db:"-"`
	IsMetaDataInited bool       `db:"-"`
}

func (p *Peer) SetAddress(addr string) {
	p.addressLen = len(addr)
	copy(p.Address[:p.addressLen], []byte(addr))
}

func (p *Peer) AddressStr() string {
	return string(p.Address[:p.addressLen])
}

func (p *Peer) PeerIDStr() string { return string(p.PeerID[:]) }
