package snet

type PeerGroup struct {
	Arr [SNetMemberCommonLimit]PeerID
	Len int
}

func (p *PeerGroup) Reset() {
	p.Len = 0
}

func (p *PeerGroup) Append(value PeerID) {
	p.Arr[p.Len] = value
	p.Len += 1
}

func (p *PeerGroup) Slice() []PeerID {
	return p.Arr[:p.Len]
}
