package snettypes

import (
	"net"
	"soloos/common/log"
	"sync"
)

type Connection struct {
	NetConn    net.Conn
	readMutex  sync.Mutex
	writeMutex sync.Mutex
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

func (p *Connection) Close(closeResonErr error) error {
	if closeResonErr != nil {
		log.Debug("connection close", closeResonErr, p.NetConn.RemoteAddr())
	}

	var err error
	err = p.NetConn.Close()
	if err != nil {
		return err
	}

	return nil
}

func (p *Connection) ReadAcquire() {
	p.readMutex.Lock()
}

func (p *Connection) ReadRelease() {
	p.readMutex.Unlock()
}

func (p *Connection) WaitReadDone() {
	p.readMutex.Lock()
	p.readMutex.Unlock()
}

func (p *Connection) WriteAcquire() {
	p.writeMutex.Lock()
}

func (p *Connection) WriteRelease() {
	p.writeMutex.Unlock()
}

func (p *Connection) WaitWriteDone() {
	p.writeMutex.Lock()
	p.writeMutex.Unlock()
}
