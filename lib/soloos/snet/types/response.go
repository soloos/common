package types

import (
	"encoding/binary"
)

const (
	ResponseHeaderBaseSize = 14
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

type ClientResponse struct {
	BodySize uint32
}
