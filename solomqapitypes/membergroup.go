package solomqapitypes

type SOLOMQMemberGroup struct {
	Arr [8]SOLOMQMember
	Len int
}

func (p *SOLOMQMemberGroup) Reset() {
	p.Len = 0
}

func (p *SOLOMQMemberGroup) Append(value SOLOMQMember) {
	p.Arr[p.Len] = value
	p.Len += 1
}

func (p *SOLOMQMemberGroup) Slice() []SOLOMQMember {
	return p.Arr[:p.Len]
}

func (p *SOLOMQMemberGroup) SetSOLOMQMembers(solomqMembers []SOLOMQMember) {
	for i := 0; i < len(p.Arr) && i < len(solomqMembers); i++ {
		p.Append(solomqMembers[i])
	}
}
