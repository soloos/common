package snet

import "soloos/snet/types"

type NetDriver interface {
	Init() error
	Call(types.PeerUintptr, string, []byte) error
}
