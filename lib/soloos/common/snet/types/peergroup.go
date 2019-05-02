package types

type PeerGroup struct {
	Arr [8]PeerUintptr
	Len int
}

func (p *PeerGroup) Reset() {
	p.Len = 0
}

func (p *PeerGroup) Append(value PeerUintptr) {
	p.Arr[p.Len] = value
	p.Len += 1
}

func (p *PeerGroup) Slice() []PeerUintptr {
	return p.Arr[:p.Len]
}
