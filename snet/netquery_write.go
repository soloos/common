package snet

import (
	"io"
	"soloos/common/iron"
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

func (p *NetQuery) WriteSNetReqHeader(reqID uint64, url string, bodySize, reqParamSize uint32) error {
	p.Conn.WriteAcquire()

	p.ReqID = reqID
	p.BodySize = bodySize
	p.ParamSize = reqParamSize
	p.ConnBytesLeft = p.BodySize

	var (
		reqHeader SNetReqHeader
		off, n    int
		err       error
	)
	reqHeader.SetID(p.ReqID)
	reqHeader.SetVersion(SNetVersion)
	reqHeader.SetBodySize(p.BodySize)
	reqHeader.SetParamSize(p.ParamSize)
	reqHeader.SetUrl(url)

	for off = 0; off < len(reqHeader.Base); off += n {
		n, err = p.Conn.NetConn.Write(reqHeader.Base[:])
		if err != nil {
			p.Conn.WriteRelease()
			return err
		}
	}

	for off = 0; off < len(reqHeader.SNetReqHeaderBodyBs); off += n {
		n, err = p.Conn.NetConn.Write(reqHeader.SNetReqHeaderBodyBs)
		if err != nil {
			p.Conn.WriteRelease()
			return err
		}
	}

	err = p.afterWriteHeader(bodySize)
	return err
}

func (p *NetQuery) WriteSNetRespHeader(reqID uint64, bodySize, reqParamSize uint32) error {
	p.Conn.WriteAcquire()

	p.ReqID = reqID
	p.BodySize = bodySize
	p.ParamSize = reqParamSize
	p.ConnBytesLeft = p.BodySize

	var (
		respHeader SNetRespHeader
		off, n     int
		err        error
	)
	respHeader.SetID(p.ReqID)
	respHeader.SetVersion(SNetVersion)
	respHeader.SetBodySize(p.BodySize)
	respHeader.SetParamSize(p.ParamSize)

	for off = 0; off < len(respHeader.Base); off += n {
		n, err = p.Conn.NetConn.Write(respHeader.Base[:])
		if err != nil {
			p.Conn.WriteRelease()
			return err
		}
	}
	for off = 0; off < len(respHeader.SNetRespHeaderBodyBs); off += n {
		n, err = p.Conn.NetConn.Write(respHeader.SNetRespHeaderBodyBs)
		if err != nil {
			p.Conn.WriteRelease()
			return err
		}
	}

	err = p.afterWriteHeader(bodySize)
	return err
}

func (p *NetQuery) SimpleResponse(reqID uint64, respBody []byte) error {
	var err error
	err = p.WriteSNetRespHeader(reqID, uint32(len(respBody)), uint32(len(respBody)))
	if err != nil {
		return err
	}

	err = p.WriteAll(respBody)
	if err != nil {
		return err
	}

	return nil
}

func (p *NetQuery) ResponseWithOffheap(reqID uint64, resp interface{}, offheapBodySize int) error {
	var err error
	var respBody = iron.MustSpecMarshalResponseErr(resp, nil)

	err = p.WriteSNetRespHeader(reqID,
		uint32(len(respBody)+offheapBodySize),
		uint32(len(respBody)))
	if err != nil {
		return err
	}

	err = p.WriteAll(respBody)
	if err != nil {
		return err
	}

	return nil
}
