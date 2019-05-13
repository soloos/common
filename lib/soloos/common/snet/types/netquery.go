package types

type NetQuery struct {
	ReqID         uint64
	BodySize      uint32
	ParamSize     uint32
	ConnBytesLeft uint32
	conn          *Connection
}

func (p *NetQuery) Init(conn *Connection) {
	p.conn = conn
}

func (p *NetQuery) ConnConnect(address string) error {
	return p.conn.Connect(address)
}

func (p *NetQuery) ConnClose(closeResonErr error) error {
	return p.conn.Close(closeResonErr)
}
