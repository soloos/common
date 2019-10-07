package snet

import (
	"io"
	"soloos/common/util"
)

func (p *NetQuery) afterReadHeader(maxMessageLength, bodySize uint32, netVersion byte) error {
	// prepare & check req
	if netVersion != SNetVersion {
		p.Conn.ReadRelease()
		return ErrWrongVersion
	}

	if bodySize > maxMessageLength {
		p.Conn.ReadRelease()
		return ErrMessageTooLong
	}

	if p.ConnBytesLeft == 0 {
		p.Conn.ReadRelease()
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
		n, err = p.Conn.NetConn.Read(b)
	} else {
		n, err = p.Conn.NetConn.Read(b[:int(p.ConnBytesLeft)])
	}

	if err != nil {
		p.Conn.ReadRelease()
		return n, err
	}

	p.ConnBytesLeft -= uint32(n)
	if p.ConnBytesLeft == 0 {
		p.Conn.ReadRelease()
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

func (p *NetQuery) ReadSNetReqHeader(maxMessageLength uint32, reqHeader *SNetReqHeader) error {
	p.Conn.ReadAcquire()

	var (
		offset, n int
		err       error
	)

	for offset = 0; offset < len(reqHeader.Base); offset += n {
		n, err = p.Conn.NetConn.Read(reqHeader.Base[offset:len(reqHeader.Base)])
		if err != nil {
			p.Conn.ReadRelease()
			return err
		}
	}

	var headerBodyBytes = make([]byte, reqHeader.HeaderBodySize())
	for offset = 0; offset < len(headerBodyBytes); offset += n {
		n, err = p.Conn.NetConn.Read(headerBodyBytes[offset:len(headerBodyBytes)])
		if err != nil {
			p.Conn.ReadRelease()
			return err
		}
	}
	reqHeader.SetHeaderBody(headerBodyBytes)

	p.ReqID = reqHeader.ID()
	p.BodySize = reqHeader.BodySize()
	p.ParamSize = reqHeader.ParamSize()
	p.ConnBytesLeft = p.BodySize

	err = p.afterReadHeader(maxMessageLength, reqHeader.BodySize(), reqHeader.Version())
	return err
}

func (p *NetQuery) ReadSNetRespHeader(maxMessageLength uint32, respHeader *SNetRespHeader) error {
	p.Conn.ReadAcquire()

	var (
		offset, n int
		err       error
	)

	for offset = 0; offset < len(respHeader.Base); offset += n {
		n, err = p.Conn.NetConn.Read(respHeader.Base[offset:len(respHeader.Base)])
		if err != nil {
			p.Conn.ReadRelease()
			return err
		}
	}

	var headerBodyBytes = make([]byte, respHeader.HeaderBodySize())
	for offset = 0; offset < len(headerBodyBytes); offset += n {
		n, err = p.Conn.NetConn.Read(headerBodyBytes[offset:len(headerBodyBytes)])
		if err != nil {
			p.Conn.ReadRelease()
			return err
		}
	}
	respHeader.SetHeaderBody(headerBodyBytes)

	p.ReqID = respHeader.ID()
	p.BodySize = respHeader.BodySize()
	p.ParamSize = respHeader.ParamSize()
	p.ConnBytesLeft = p.BodySize

	err = p.afterReadHeader(maxMessageLength, respHeader.BodySize(), respHeader.Version())
	return err
}

func (p *NetQuery) EnsureServiceReadDone() {
	if p.ConnBytesLeft > 0 {
		p.SkipReadRemaining()
	}
}
