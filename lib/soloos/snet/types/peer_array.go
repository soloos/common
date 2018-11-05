package types

type PeerUintptrArray8 struct {
	Arr [8]PeerUintptr
	Len int
}

func (p *PeerUintptrArray8) Append(value PeerUintptr) {
	p.Arr[p.Len] = value
	p.Len += 1
}
