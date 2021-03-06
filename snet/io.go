package snet

import (
	"unsafe"
)

type OffheapFastCopyer struct {
	OffheapBytes uintptr
	CopyOffset   int
	CopyEnd      int
}

func (p *OffheapFastCopyer) BodySize() int {
	return p.CopyEnd - p.CopyOffset
}

func (p *OffheapFastCopyer) Copy(netQuery *NetQuery) error {
	if p.OffheapBytes == 0 {
		return nil
	}

	if p.CopyOffset >= p.CopyEnd {
		return nil
	}

	var data = (*((*[1 << 31]byte)(unsafe.Pointer(p.OffheapBytes))))[p.CopyOffset:p.CopyEnd]
	return netQuery.WriteAll(data)
}
