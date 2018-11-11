package types

import (
	"io"
)

func (p *Connection) WriteAcquire() {
	p.writeMutex.Lock()
}

func (p *Connection) WriteRelease() {
	p.writeMutex.Unlock()
}

func (p *Connection) afterWriteHeader(contentLen uint32) error {
	p.LastWriteLimit = contentLen
	if p.LastWriteLimit == 0 {
		p.WriteRelease()
	}
	return nil
}

func (p *Connection) WriteRequestHeader(requestID uint64, serviceID string, contentLen uint32) error {
	p.WriteAcquire()

	var (
		header RequestHeader
		off, n int
		err    error
	)
	header.SetID(requestID)
	header.SetVersion(SNetVersion)
	header.SetContentLen(contentLen)
	header.SetServiceID(serviceID)
	for off = 0; off < len(header); off += n {
		n, err = p.NetConn.Write(header[:])
		if err != nil {
			p.WriteRelease()
			return err
		}
	}

	return p.afterWriteHeader(contentLen)
}

func (p *Connection) WriteResponseHeader(requestID uint64, contentLen uint32) error {
	p.WriteAcquire()

	var (
		header ResponseHeader
		off, n int
		err    error
	)
	header.SetID(requestID)
	header.SetVersion(SNetVersion)
	header.SetContentLen(contentLen)
	for off = 0; off < len(header); off += n {
		n, err = p.NetConn.Write(header[:])
		if err != nil {
			p.WriteRelease()
			return err
		}
	}

	return p.afterWriteHeader(contentLen)
}

func (p *Connection) Write(b []byte) (int, error) {
	if p.LastWriteLimit == 0 {
		return 0, io.EOF
	}

	if len(b) == 0 {
		return 0, nil
	}

	var (
		n   int
		err error
	)

	if p.LastWriteLimit > uint32(len(b)) {
		n, err = p.NetConn.Write(b)
	} else {
		n, err = p.NetConn.Write(b[:int(p.LastWriteLimit)])
	}

	if err != nil {
		p.WriteRelease()
		return n, err
	}

	p.LastWriteLimit -= uint32(n)
	if p.LastWriteLimit == 0 {
		p.WriteRelease()
	}

	return n, nil
}

func (p *Connection) WriteAll(b []byte) error {
	var (
		n   int
		err error
	)

	for off := 0; off < len(b); off += n {
		n, err = p.Write(b[off:])
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
	}

	return nil
}

func (p *Connection) Response(requestID uint64, resp []byte) error {
	var err error
	err = p.WriteResponseHeader(requestID, uint32(len(resp)))
	if err != nil {
		return err
	}

	err = p.WriteAll(resp)
	if err != nil {
		return err
	}

	return nil
}
