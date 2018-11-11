package types

import (
	"reflect"
	"unsafe"
)

type OffheapFastCopyer struct {
	OffheapBytes reflect.SliceHeader
	CopyOffset   int
	CopyEnd      int
}

func (p *OffheapFastCopyer) ContentLen() int {
	return p.CopyEnd - p.CopyOffset
}

func (p *OffheapFastCopyer) Copy(conn *Connection) error {
	if p.OffheapBytes.Data == 0 {
		return nil
	}

	if p.CopyOffset >= p.CopyEnd {
		return nil
	}

	return conn.WriteAll((*((*[]byte)(unsafe.Pointer(&p.OffheapBytes))))[p.CopyOffset:p.CopyEnd])
}
