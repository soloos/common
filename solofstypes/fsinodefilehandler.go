package solofstypes

type FsINodeFileHandlerID = uint64

type FsINodeFileHandler struct {
	FsINodeIno      FsINodeIno
	AppendPosition uint64
	ReadPosition   uint64
}

func (p *FsINodeFileHandler) Reset() {
	p.AppendPosition = 0
	p.ReadPosition = 0
}
