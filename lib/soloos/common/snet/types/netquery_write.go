package types

import (
	"io"
)

func (p *NetQuery) afterWriteHeader(bodySize uint32) error {
	if p.ConnBytesLeft == 0 {
		p.conn.WriteRelease()
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
		n, err = p.conn.NetConn.Write(b)
	} else {
		n, err = p.conn.NetConn.Write(b[:int(p.ConnBytesLeft)])
	}

	if err != nil {
		p.conn.WriteRelease()
		return n, err
	}

	p.ConnBytesLeft -= uint32(n)
	if p.ConnBytesLeft == 0 {
		p.conn.WriteRelease()
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

func (p *NetQuery) WriteRequestHeader(reqID uint64, serviceID string, bodySize, reqParamSize uint32) error {
	p.conn.WriteAcquire()

	p.ReqID = reqID
	p.BodySize = bodySize
	p.ParamSize = reqParamSize
	p.ConnBytesLeft = p.BodySize

	var (
		reqHeader RequestHeader
		off, n    int
		err       error
	)
	reqHeader.SetID(p.ReqID)
	reqHeader.SetVersion(SNetVersion)
	reqHeader.SetBodySize(p.BodySize)
	reqHeader.SetParamSize(p.ParamSize)
	reqHeader.SetServiceID(serviceID)
	for off = 0; off < len(reqHeader); off += n {
		n, err = p.conn.NetConn.Write(reqHeader[:])
		if err != nil {
			p.conn.WriteRelease()
			return err
		}
	}

	return p.afterWriteHeader(bodySize)
}

func (p *NetQuery) WriteResponseHeader(reqID uint64, bodySize, reqParamSize uint32) error {
	p.conn.WriteAcquire()

	p.ReqID = reqID
	p.BodySize = bodySize
	p.ParamSize = reqParamSize
	p.ConnBytesLeft = p.BodySize

	var (
		respHeader ResponseHeader
		off, n     int
		err        error
	)
	respHeader.SetID(p.ReqID)
	respHeader.SetVersion(SNetVersion)
	respHeader.SetBodySize(p.BodySize)
	respHeader.SetParamSize(p.ParamSize)
	for off = 0; off < len(respHeader); off += n {
		n, err = p.conn.NetConn.Write(respHeader[:])
		if err != nil {
			p.conn.WriteRelease()
			return err
		}
	}

	return p.afterWriteHeader(bodySize)
}

func (p *NetQuery) SimpleResponse(reqID uint64, respBody []byte) error {
	var err error
	err = p.WriteResponseHeader(reqID, uint32(len(respBody)), uint32(len(respBody)))
	if err != nil {
		return err
	}

	err = p.WriteAll(respBody)
	if err != nil {
		return err
	}

	return nil
}

func (p *NetQuery) ResponseHeaderParam(reqID uint64, param []byte, offheapBodySize int) error {
	var err error
	err = p.WriteResponseHeader(reqID,
		uint32(len(param)+offheapBodySize),
		uint32(len(param)))
	if err != nil {
		return err
	}

	err = p.WriteAll(param)
	if err != nil {
		return err
	}

	return nil
}

func (p *NetQuery) Response(reqID uint64, param []byte, offheapBody []byte) error {
	var err error
	err = p.WriteResponseHeader(reqID,
		uint32(len(param)+len(offheapBody)),
		uint32(len(param)))
	if err != nil {
		return err
	}

	err = p.WriteAll(param)
	if err != nil {
		return err
	}

	err = p.WriteAll(offheapBody)
	if err != nil {
		return err
	}

	return nil
}
