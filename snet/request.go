package snet

import (
	"encoding/binary"
	"encoding/json"
	"soloos/common/util"
	"unsafe"
)

const (
	SNetReqHeaderBaseBaseSize = 24
	SNetReqHeaderBaseSize     = uint32(unsafe.Sizeof(SNetReqHeaderBase{}))
)

type SNetReqHeaderBase [SNetReqHeaderBaseBaseSize]byte
type SNetReqHeaderBody struct {
	Url string
}

type SNetReqHeader struct {
	Base SNetReqHeaderBase
	SNetReqHeaderBody
	SNetReqHeaderBodyBs []byte
}

func (p *SNetReqHeader) SetVersion(version byte) {
	p.Base[0] = version
}

func (p *SNetReqHeader) Version() byte {
	return p.Base[0]
}

func (p *SNetReqHeader) IsHeartbeat() bool {
	return p.Base[1]&0x01 == 0x01
}

func (p *SNetReqHeader) SetHeartbeat(isHeartbeat bool) {
	if isHeartbeat {
		p.Base[1] = p.Base[1] | 0x01
	} else {
		p.Base[1] = p.Base[1] &^ 0x01
	}
}

func (p *SNetReqHeader) IsNotExceptResponse() bool {
	return p.Base[1]&0x02 == 0x02
}

func (p *SNetReqHeader) SetNotExceptResponse(isNotExceptResponse bool) {
	if isNotExceptResponse {
		p.Base[1] = p.Base[1] | 0x02
	} else {
		p.Base[1] = p.Base[1] &^ 0x02
	}
}

func (p *SNetReqHeader) ID() uint64 {
	return binary.BigEndian.Uint64(p.Base[2:10])
}

func (p *SNetReqHeader) SetID(seq uint64) {
	binary.BigEndian.PutUint64(p.Base[2:10], seq)
}

func (p *SNetReqHeader) HeaderBodySize() uint32 {
	return binary.BigEndian.Uint32(p.Base[10:14])
}

func (p *SNetReqHeader) SetHeaderBodySize(size uint32) {
	binary.BigEndian.PutUint32(p.Base[10:14], size)
}

func (p *SNetReqHeader) ParamSize() uint32 {
	return binary.BigEndian.Uint32(p.Base[14:18])
}

func (p *SNetReqHeader) SetParamSize(reqParamSize uint32) {
	binary.BigEndian.PutUint32(p.Base[14:18], reqParamSize)
}

func (p *SNetReqHeader) BodySize() uint32 {
	return binary.BigEndian.Uint32(p.Base[18:24])
}

func (p *SNetReqHeader) SetBodySize(bodySize uint32) {
	binary.BigEndian.PutUint32(p.Base[18:24], bodySize)
}

func (p *SNetReqHeader) SetHeaderBody(headerBodyBs []byte) {
	p.SNetReqHeaderBodyBs = headerBodyBs
	json.Unmarshal(p.SNetReqHeaderBodyBs, &p.SNetReqHeaderBody)
}

func (p *SNetReqHeader) SetUrl(url string) {
	p.SNetReqHeaderBody.Url = url
	var err error
	p.SNetReqHeaderBodyBs, err = json.Marshal(p.SNetReqHeaderBody)
	util.AssertErrIsNil(err)
	p.SetHeaderBodySize(uint32(len(p.SNetReqHeaderBodyBs)))
}

type SNetReqContext struct {
	IsResponseInService bool
	Header              SNetReqHeader
	NetQuery
}

func (p *SNetReqContext) SetResponseInService() {
	p.IsResponseInService = true
}

func (p *SNetReqContext) ReadHeader(maxMessageLength uint32) error {
	return p.NetQuery.ReadSNetReqHeader(maxMessageLength, &p.Header)
}

type SNetReq struct {
	NetQuery
	Url         string
	Param       []byte
	OffheapBody OffheapFastCopyer
}

func (p *SNetReq) Init(reqID uint64, conn *Connection, url string) {
	p.NetQuery.Init(conn)
	p.NetQuery.ReqID = reqID
	p.Url = url
}
