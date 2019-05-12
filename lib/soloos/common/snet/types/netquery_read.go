package types

import (
	"io"
	"soloos/common/util"
)

func (p *NetQuery) afterReadHeader(maxMessageLength, bodySize uint32, netVersion byte) error {
	// prepare & check req
	if netVersion != SNetVersion {
		return ErrWrongVersion
	}

	p.Conn.LastReadLimit = bodySize
	if bodySize > maxMessageLength {
		return ErrMessageTooLong
	}

	return nil
}

func (p *NetQuery) Read(b []byte) (int, error) {
	if p.ConnBytesLimit == 0 {
		return 0, io.EOF
	}

	var (
		n   int
		err error
	)

	if p.ConnBytesLimit > uint32(len(b)) {
		n, err = p.Conn.NetConn.Read(b)
	} else {
		n, err = p.Conn.NetConn.Read(b[:int(p.ConnBytesLimit)])
	}

	if err != nil {
		p.Conn.ReadRelease()
		return n, err
	}

	p.ConnBytesLimit -= uint32(n)
	if p.ConnBytesLimit == 0 {
		p.Conn.ReadRelease()
	}

	return n, err
}

func (p *NetQuery) SkipReadRemaining() error {
	var err error
	for p.ConnBytesLimit > 0 {
		if p.ConnBytesLimit > uint32(len(util.DevNullBuf)) {
			err = p.ReadAll(util.DevNullBuf[:])
		} else {
			err = p.ReadAll(util.DevNullBuf[:p.ConnBytesLimit])
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *NetQuery) ReadAll(b []byte) error {
	var (
		off, n int
		err    error
	)

	for off = 0; off < len(b); off += n {
		n, err = p.Read(b[off:])
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *NetQuery) ReadRequestHeader(maxMessageLength uint32, header *RequestHeader) error {
	p.Conn.ReadAcquire()

	var (
		offset, n int
		err       error
	)

	for offset = 0; offset < len(header); offset += n {
		n, err = p.Conn.NetConn.Read(header[offset:len(header)])
		if err != nil {
			p.Conn.ReadRelease()
			return err
		}
	}

	return p.afterReadHeader(maxMessageLength, header.BodySize(), header.Version())
}

func (p *NetQuery) ReadResponseHeader(maxMessageLength uint32, header *ResponseHeader) error {
	p.Conn.ReadAcquire()

	var (
		offset, n int
		err       error
	)

	for offset = 0; offset < len(header); offset += n {
		n, err = p.Conn.NetConn.Read(header[offset:len(header)])
		if err != nil {
			p.Conn.ReadRelease()
			return err
		}
	}

	return p.afterReadHeader(maxMessageLength, header.BodySize(), header.Version())
}
