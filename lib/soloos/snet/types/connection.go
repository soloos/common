package types

import (
	"io"
	"net"
	"sync"
)

type Connection struct {
	netConn              net.Conn
	ContinueReadSig      sync.WaitGroup
	LastRequestReadLimit uint32
	LastRequestError     error

	readMutex  sync.Mutex
	writeMutex sync.Mutex
}

func (p *Connection) Init(netConn net.Conn) {
	p.netConn = netConn
}

func (p *Connection) ReadAcquire() {
	p.readMutex.Lock()
}

func (p *Connection) ReadRelease() {
	p.readMutex.Unlock()
}

func (p *Connection) ReadRequestHeader(header *RequestHeader) error {
	var (
		offset, n int
		err       error
	)

	for offset = 0; offset < len(header); offset += n {
		n, err = p.netConn.Read(header[offset:len(header)])
		if err != nil && err != io.EOF {
			return err
		}
	}
	return nil
}

func (p *Connection) ReadResponseHeader(header *ResponseHeader) error {
	var (
		offset, n int
		err       error
	)

	for offset = 0; offset < len(header); offset += n {
		n, err = p.netConn.Read(header[offset:len(header)])
		if err != nil && err != io.EOF {
			return err
		}
	}
	return nil
}

func (p *Connection) Read(b []byte) (int, error) {
	var n int
	if p.LastRequestReadLimit > uint32(len(b)) {
		n, p.LastRequestError = p.netConn.Read(b)
	} else {
		n, p.LastRequestError = p.netConn.Read(b[:int(p.LastRequestReadLimit)])
	}

	if p.LastRequestError != nil {
		p.ContinueReadSig.Done()
		return n, p.LastRequestError
	}

	p.LastRequestReadLimit -= uint32(n)
	if p.LastRequestReadLimit == 0 {
		p.ContinueReadSig.Done()
	}
	return n, p.LastRequestError
}

func (p *Connection) ReadAll(b []byte) error {
	var (
		off, n int
		err    error
	)
	for off = 0; off < len(b); off += n {
		n, err = p.Read(b[off:])
		if err != nil && err != io.EOF {
			return err
		}
	}
	p.ContinueReadSig.Wait()
	return nil
}

func (p *Connection) SkipReadAllRest(b *[512]byte) error {
	var (
		n   int
		err error
	)
	for p.LastRequestReadLimit > 0 {
		if p.LastRequestReadLimit > uint32(len(b)) {
			n, err = p.netConn.Read(b[:])
		} else {
			n, err = p.netConn.Read(b[:p.LastRequestReadLimit])
		}
		if err != nil && err != io.EOF {
			return err
		}
		p.LastRequestReadLimit -= uint32(n)
	}
	return nil
}

func (p *Connection) WriteAcquire() {
	p.writeMutex.Lock()
}

func (p *Connection) WriteRelease() {
	p.writeMutex.Unlock()
}

func (p *Connection) WriteHeader(requestID uint64, contentLen uint32) error {
	var (
		header ResponseHeader
		err    error
	)
	header.SetID(requestID)
	header.SetVersion(SNetVersion)
	header.SetContentLen(contentLen)
	err = p.WriteAll(header[:])
	if err != nil {
		return err
	}

	return nil
}

func (p *Connection) WriteAll(b []byte) error {
	var (
		n   int
		err error
	)

	for off := 0; off < len(b); off += n {
		n, err = p.netConn.Write(b)
		if err != nil && err != io.EOF {
			return err
		}
	}

	return nil
}

func (p *Connection) Close() error {
	return p.netConn.Close()
}
