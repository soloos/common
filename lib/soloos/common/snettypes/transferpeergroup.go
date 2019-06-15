package snettypes

type TransferPeer struct {
	PeerID        PeerID
	TransferCount byte
}

type TransferPeerGroup struct {
	Arr [8]TransferPeer
	Len int
}

func (p *TransferPeerGroup) Reset() {
	p.Len = 0
}

func (p *TransferPeerGroup) Append(value PeerID, count int) {
	p.Arr[p.Len].PeerID = value
	p.Arr[p.Len].TransferCount = byte(count)
	p.Len += 1
}

func (p *TransferPeerGroup) Slice() []TransferPeer {
	return p.Arr[:p.Len]
}
