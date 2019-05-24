package snettypes

import (
	"encoding/binary"
	"unsafe"
)

const (
	RequestHeaderBaseSize = 18
	RequestHeaderSize     = uint32(unsafe.Sizeof(RequestHeader{}))
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

func (p *RequestHeader) BodySize() uint32 {
	return binary.BigEndian.Uint32(p[10:14])
}

func (p *RequestHeader) SetBodySize(bodySize uint32) {
	binary.BigEndian.PutUint32(p[10:14], bodySize)
}

func (p *RequestHeader) ParamSize() uint32 {
	return binary.BigEndian.Uint32(p[14:18])
}

func (p *RequestHeader) SetParamSize(reqParamSize uint32) {
	binary.BigEndian.PutUint32(p[14:18], reqParamSize)
}

func (p *RequestHeader) ServiceID(ret *ServiceID) {
	copy((*ret)[:], p[18:18+ServiceIDLen])
	return
}

func (p *RequestHeader) SetServiceID(serviceID string) {
	copy(p[18:18+ServiceIDLen], []byte(serviceID))
}

type Request struct {
	NetQuery
	ServiceID   string
	Param       []byte
	OffheapBody OffheapFastCopyer
}

func (p *Request) Init(reqID uint64, conn *Connection, serviceID string) {
	p.NetQuery.ReqID = reqID
	p.NetQuery.Init(conn)
	p.ServiceID = serviceID
}
