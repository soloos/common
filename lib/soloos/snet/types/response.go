package types

import (
	"encoding/binary"
	"soloos/util/offheap"
	"unsafe"
)

const (
	ResponseHeaderBaseSize = 14
	ResponseHeaderSize     = uint32(unsafe.Sizeof(ResponseHeader{}))
)

type ResponseHeader [ResponseHeaderBaseSize]byte

func (p *ResponseHeader) SetVersion(version byte) {
	p[0] = version
}

func (p *ResponseHeader) Version() byte {
	return p[0]
}

func (p *ResponseHeader) ID() uint64 {
	return binary.BigEndian.Uint64(p[2:10])
}

func (p *ResponseHeader) SetID(seq uint64) {
	binary.BigEndian.PutUint64(p[2:10], seq)
}

func (p *ResponseHeader) ContentLen() uint32 {
	return binary.BigEndian.Uint32(p[10:14])
}

func (p *ResponseHeader) SetContentLen(contentLen uint32) {
	binary.BigEndian.PutUint32(p[10:14], contentLen)
}

type Response struct {
	BodySize       uint32
	NetConnReadSig offheap.MutexUintptr
}
