package solomqapitypes

import "soloos/common/snettypes"

type SolomqMemberGroup struct {
	Arr [snettypes.SNetMemberCommonLimit]SolomqMember
	Len int
}

func (p *SolomqMemberGroup) Reset() {
	p.Len = 0
}

func (p *SolomqMemberGroup) Append(value SolomqMember) {
	p.Arr[p.Len] = value
	p.Len += 1
}

func (p *SolomqMemberGroup) Slice() []SolomqMember {
	return p.Arr[:p.Len]
}

func (p *SolomqMemberGroup) SetSolomqMembers(solomqMembers []SolomqMember) {
	for i := 0; i < len(p.Arr) && i < len(solomqMembers); i++ {
		p.Append(solomqMembers[i])
	}
}
