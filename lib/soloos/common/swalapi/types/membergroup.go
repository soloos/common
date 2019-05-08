package types

type SWALMemberGroup struct {
	Arr [8]SWALMember
	Len int
}

func (p *SWALMemberGroup) Reset() {
	p.Len = 0
}

func (p *SWALMemberGroup) Append(value SWALMember) {
	p.Arr[p.Len] = value
	p.Len += 1
}

func (p *SWALMemberGroup) Slice() []SWALMember {
	return p.Arr[:p.Len]
}

func (p *SWALMemberGroup) SetSWALMembers(swalMembers []SWALMember) {
	for i := 0; i < len(p.Arr) && i < len(swalMembers); i++ {
		p.Append(swalMembers[i])
	}
}
