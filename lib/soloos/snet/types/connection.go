package types

import (
	"net"
	"sync"
)

type Connection struct {
	NetConn net.Conn

	ContinueReadSig sync.Mutex

	LastReadLimit uint32
	readMutex     sync.Mutex

	LastWriteLimit uint32
	writeMutex     sync.Mutex
}

func (p *Connection) SetNetConn(netConn net.Conn) {
	p.NetConn = netConn
}

func (p *Connection) Connect(address string) error {
	var err error
	if p.NetConn != nil {
		err = p.NetConn.Close()
		if err != nil {
			return err
		}
	}

	p.NetConn, err = net.Dial("tcp", address)
	if err != nil {
		return err
	}

	return nil
}

func (p *Connection) Close() error {
	var err error
	err = p.NetConn.Close()
	if err != nil {
		return err
	}

	return nil
}
