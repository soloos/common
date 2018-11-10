package srpc

import (
	"net"
	"soloos/snet/types"
)

type Client struct {
	Conn types.Connection
}

func (p *Client) Init(address string) error {
	var (
		err     error
		netConn net.Conn
	)

	netConn, err = net.Dial("tcp", address)
	if err != nil {
		return err
	}

	p.Conn.Init(netConn)

	return nil
}
