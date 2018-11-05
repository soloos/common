package types

import "unsafe"

const (
	PeerStructSize = unsafe.Sizeof(Peer{})
)

type PeerID = [64]byte

type PeerUintptr uintptr

func (u PeerUintptr) Ptr() *Peer { return (*Peer)(unsafe.Pointer(u)) }

type Peer struct {
	ID              PeerID
	Address         [128]byte
	ServiceProtocol [32]byte
}

func InitPeerID(peerID *PeerID) {
}
