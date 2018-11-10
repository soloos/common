package types

import (
	"io"
	"reflect"
	"unsafe"
)

type OffheapFastCopyer struct {
	OffheapBytes reflect.SliceHeader
	CopyOffset   int
	CopyEnd      int
}

func (p *OffheapFastCopyer) Copy(conn *Connection) error {
	if p.OffheapBytes.Data == 0 {
		return io.EOF
	}

	if p.CopyOffset >= p.CopyEnd {
		return io.EOF
	}
	return conn.WriteAll((*((*[]byte)(unsafe.Pointer(&p.OffheapBytes))))[p.CopyOffset:p.CopyEnd])
}
