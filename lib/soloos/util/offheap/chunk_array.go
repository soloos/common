package offheap

type ChunkUintptrArray8 struct {
	Arr [8]ChunkUintptr
	Len int
}

func (p *ChunkUintptrArray8) Append(value ChunkUintptr) {
	p.Arr[p.Len] = value
	p.Len += 1
}

type ChunkUintptrArray16 struct {
	Arr [16]ChunkUintptr
	Len int
}

func (p *ChunkUintptrArray16) Append(value ChunkUintptr) {
	p.Arr[p.Len] = value
	p.Len += 1
}

type ChunkUintptrArray32 struct {
	Arr [32]ChunkUintptr
	Len int
}

func (p *ChunkUintptrArray32) Append(value ChunkUintptr) {
	p.Arr[p.Len] = value
	p.Len += 1
}

type ChunkUintptrArray64 struct {
	Arr [64]ChunkUintptr
	Len int
}

func (p *ChunkUintptrArray64) Append(value ChunkUintptr) {
	p.Arr[p.Len] = value
	p.Len += 1
}
