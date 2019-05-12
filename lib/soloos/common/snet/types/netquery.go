package types

type NetQuery struct {
	ReqID          uint64
	ReqBodySize    uint32
	ReqParamSize   uint32
	Conn           *Connection
	ConnBytesLimit uint32
}
