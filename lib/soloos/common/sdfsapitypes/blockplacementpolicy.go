package sdfsapitypes

const (
	BlockPlacementPolicyDefault = byte(iota)
	BlockPlacementPolicySWAL
	MemBlockPlacementPolicyHeaderBytesNum = 2
)

type MemBlockPlacementPolicy [16]byte

func (p *MemBlockPlacementPolicy) SetType(t byte) {
	(*p)[0] = t
}

func (p *MemBlockPlacementPolicy) GetType() byte {
	return p[0]
}
