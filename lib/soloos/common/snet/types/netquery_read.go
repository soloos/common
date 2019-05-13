package types

import (
	"io"
	"soloos/common/util"
)

func (p *NetQuery) afterReadHeader(maxMessageLength, bodySize uint32, netVersion byte) error {
	// prepare & check req
	if netVersion != SNetVersion {
		p.conn.ReadRelease()
		return ErrWrongVersion
	}

	if bodySize > maxMessageLength {
		p.conn.ReadRelease()
		return ErrMessageTooLong
	}

	if p.ConnBytesLeft == 0 {
		p.conn.ReadRelease()
	}

	return nil
}

func (p *NetQuery) Read(b []byte) (int, error) {
	if p.ConnBytesLeft == 0 {
		return 0, io.EOF
	}

	var (
		n   int
		err error
	)

	if p.ConnBytesLeft > uint32(len(b)) {
		n, err = p.conn.NetConn.Read(b)
	} else {
		n, err = p.conn.NetConn.Read(b[:int(p.ConnBytesLeft)])
	}

	if err != nil {
		p.conn.ReadRelease()
		return n, err
	}

	p.ConnBytesLeft -= uint32(n)
	if p.ConnBytesLeft == 0 {
		p.conn.ReadRelease()
	}

	return n, err
}

func (p *NetQuery) SkipReadRemaining() error {
	var err error
	for p.ConnBytesLeft > 0 {
		if p.ConnBytesLeft > uint32(len(util.DevNullBuf)) {
			err = p.ReadAll(util.DevNullBuf[:])
		} else {
			err = p.ReadAll(util.DevNullBuf[:p.ConnBytesLeft])
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

func (p *NetQuery) ReadRequestHeader(maxMessageLength uint32, reqHeader *RequestHeader) error {
	p.conn.ReadAcquire()

	var (
		offset, n int
		err       error
	)

	for offset = 0; offset < len(reqHeader); offset += n {
		n, err = p.conn.NetConn.Read(reqHeader[offset:len(reqHeader)])
		if err != nil {
			p.conn.ReadRelease()
			return err
		}
	}

	p.ReqID = reqHeader.ID()
	p.BodySize = reqHeader.BodySize()
	p.ParamSize = reqHeader.ParamSize()
	p.ConnBytesLeft = p.BodySize

	return p.afterReadHeader(maxMessageLength, reqHeader.BodySize(), reqHeader.Version())
}

func (p *NetQuery) ReadResponseHeader(maxMessageLength uint32, respHeader *ResponseHeader) error {
	p.conn.ReadAcquire()

	var (
		offset, n int
		err       error
	)

	for offset = 0; offset < len(respHeader); offset += n {
		n, err = p.conn.NetConn.Read(respHeader[offset:len(respHeader)])
		if err != nil {
			p.conn.ReadRelease()
			return err
		}
	}

	p.ReqID = respHeader.ID()
	p.BodySize = respHeader.BodySize()
	p.ParamSize = respHeader.ParamSize()
	p.ConnBytesLeft = p.BodySize

	return p.afterReadHeader(maxMessageLength, respHeader.BodySize(), respHeader.Version())
}

func (p *NetQuery) EnsureServiceReadDone() {
	if p.ConnBytesLeft > 0 {
		p.SkipReadRemaining()
	}
}
