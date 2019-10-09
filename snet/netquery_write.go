package snet

import (
	"io"
)

func (p *NetQuery) afterWriteHeader(bodySize uint32) error {
	if p.ConnBytesLeft == 0 {
		p.Conn.WriteRelease()
	}
	return nil
}

func (p *NetQuery) Write(b []byte) (int, error) {
	if p.ConnBytesLeft == 0 {
		return 0, io.EOF
	}

	if len(b) == 0 {
		return 0, nil
	}

	var (
		n   int
		err error
	)

	if p.ConnBytesLeft > uint32(len(b)) {
		n, err = p.Conn.NetConn.Write(b)
	} else {
		n, err = p.Conn.NetConn.Write(b[:int(p.ConnBytesLeft)])
	}

	if err != nil {
		p.Conn.WriteRelease()
		return n, err
	}

	p.ConnBytesLeft -= uint32(n)
	if p.ConnBytesLeft == 0 {
		p.Conn.WriteRelease()
	}

	return n, nil
}

func (p *NetQuery) WriteAll(b []byte) error {
	var (
		n   int
		err error
	)

	for off := 0; off < len(b); off += n {
		n, err = p.Write(b[off:])
		if err != nil {
			return err
		}
	}

	return nil
}
