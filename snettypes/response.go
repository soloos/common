package snettypes

import (
	"encoding/binary"
	"encoding/json"
	"soloos/common/iron"
	"soloos/solodb/offheap"
	"unsafe"
)

const (
	SNetRespHeaderBaseBaseSize = 24
	SNetRespHeaderBaseSize     = uint32(unsafe.Sizeof(SNetRespHeaderBase{}))
)

type SNetRespHeaderBase [SNetRespHeaderBaseBaseSize]byte
type SNetRespHeaderBody struct {
}

type SNetRespHeader struct {
	Base SNetRespHeaderBase
	SNetRespHeaderBody
	SNetRespHeaderBodyBs []byte
}

func (p *SNetRespHeader) SetVersion(version byte) {
	p.Base[0] = version
}

func (p *SNetRespHeader) Version() byte {
	return p.Base[0]
}

func (p *SNetRespHeader) ID() uint64 {
	return binary.BigEndian.Uint64(p.Base[2:10])
}

func (p *SNetRespHeader) SetID(seq uint64) {
	binary.BigEndian.PutUint64(p.Base[2:10], seq)
}

func (p *SNetRespHeader) HeaderBodySize() uint32 {
	return binary.BigEndian.Uint32(p.Base[10:14])
}

func (p *SNetRespHeader) SetHeaderBodySize(size uint32) {
	binary.BigEndian.PutUint32(p.Base[10:14], size)
}

func (p *SNetRespHeader) ParamSize() uint32 {
	return binary.BigEndian.Uint32(p.Base[14:18])
}

func (p *SNetRespHeader) SetParamSize(reqParamSize uint32) {
	binary.BigEndian.PutUint32(p.Base[14:18], reqParamSize)
}

func (p *SNetRespHeader) BodySize() uint32 {
	return binary.BigEndian.Uint32(p.Base[18:24])
}

func (p *SNetRespHeader) SetBodySize(bodySize uint32) {
	binary.BigEndian.PutUint32(p.Base[18:24], bodySize)
}

func (p *SNetRespHeader) SetHeaderBody(headerBodyBs []byte) {
	p.SNetRespHeaderBodyBs = headerBodyBs
	json.Unmarshal(p.SNetRespHeaderBodyBs, &p.SNetRespHeaderBody)
}

type SNetResp struct {
	NetQuery
	NetConnReadSig offheap.MutexUintptr
}

type IRespData = iron.IRespData
type RespDataCommon = iron.RespDataCommon
type Response = iron.Response
