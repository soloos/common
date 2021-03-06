package solofstypes

const (
	BlockPlacementPolicyDefault = byte(iota)
	BlockPlacementPolicySolomq

	MemBlockPlacementPolicyHeaderBytesNum = 2
	MemBlockPlacementPolicyBodyOff        = MemBlockPlacementPolicyHeaderBytesNum
)

type MemBlockPlacementPolicy [16]byte

func (p *MemBlockPlacementPolicy) SetType(t byte) {
	(*p)[0] = t
}

func (p *MemBlockPlacementPolicy) GetType() byte {
	return p[0]
}
