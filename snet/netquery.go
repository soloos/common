package snet

type NetQuery struct {
	ReqID         uint64
	BodySize      uint32
	ParamSize     uint32
	ConnBytesLeft uint32
	Conn          *Connection
}

func (p *NetQuery) Init(conn *Connection) {
	p.Conn = conn
}

func (p *NetQuery) ConnConnect(address string) error {
	return p.Conn.Connect(address)
}

func (p *NetQuery) ConnClose(closeResonErr error) error {
	return p.Conn.Close(closeResonErr)
}
