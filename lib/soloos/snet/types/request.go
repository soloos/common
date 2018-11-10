package types

import (
	"encoding/binary"
)

const (
	RequestHeaderBaseSize = 14
)

type RequestHeader [RequestHeaderBaseSize + ServiceIDLen]byte

func (p *RequestHeader) SetVersion(version byte) {
	p[0] = version
}

func (p *RequestHeader) Version() byte {
	return p[0]
}

func (p *RequestHeader) IsHeartbeat() bool {
	return p[1]&0x01 == 0x01
}

func (p *RequestHeader) SetHeartbeat(isHeartbeat bool) {
	if isHeartbeat {
		p[1] = p[1] | 0x01
	} else {
		p[1] = p[1] &^ 0x01
	}
}

func (p *RequestHeader) IsNotExceptResponse() bool {
	return p[1]&0x02 == 0x02
}

func (p *RequestHeader) SetNotExceptResponse(isNotExceptResponse bool) {
	if isNotExceptResponse {
		p[1] = p[1] | 0x02
	} else {
		p[1] = p[1] &^ 0x02
	}
}

func (p *RequestHeader) ID() uint64 {
	return binary.BigEndian.Uint64(p[2:10])
}

func (p *RequestHeader) SetID(seq uint64) {
	binary.BigEndian.PutUint64(p[2:10], seq)
}

func (p *RequestHeader) ContentLen() uint32 {
	return binary.BigEndian.Uint32(p[10:14])
}

func (p *RequestHeader) SetContentLen(contentLen uint32) {
	binary.BigEndian.PutUint32(p[10:14], contentLen)
}

func (p *RequestHeader) ServiceID() (ret ServiceID) {
	copy(ret[:], p[14:14+ServiceIDLen])
	return
}

func (p *RequestHeader) SetServiceID(serviceID string) {
	copy(p[14:14+ServiceIDLen], []byte(serviceID))
}

type ClientRequest struct {
	Body        []byte
	OffheapBody OffheapFastCopyer
}
