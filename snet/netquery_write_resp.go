package snet

import "soloos/common/util"

func (p *NetQuery) writeSNetRespHeaderWithLock(reqID uint64, bodySize, reqParamSize uint32) error {
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

func (p *NetQuery) SimpleResponse(reqID uint64, resp IResponse) error {
	// func (p *NetQuery) SimpleResponse(reqID uint64, resp *Response) error {
	p.Conn.WriteAcquire()

	var buf = snetCodecBytesPool.Get().(util.Buffer)
	defer snetCodecBytesPoolPut(buf)

	var err error
	err = p.Conn.Marshal(&buf, resp)
	if err != nil {
		return err
	}

	var respBody = buf.Bytes()
	err = p.writeSNetRespHeaderWithLock(reqID, uint32(len(respBody)), uint32(len(respBody)))
	if err != nil {
		return err
	}

	err = p.WriteAll(respBody)
	if err != nil {
		return err
	}

	return nil
}

func (p *NetQuery) ResponseWithOffheap(reqID uint64,
	resp IResponse,
	// resp *Response,
	offheapBodySize int) error {
	p.Conn.WriteAcquire()

	var buf = snetCodecBytesPool.Get().(util.Buffer)
	defer snetCodecBytesPoolPut(buf)

	var err error
	err = p.Conn.Marshal(&buf, resp)
	if err != nil {
		return err
	}

	var respBody = buf.Bytes()
	err = p.writeSNetRespHeaderWithLock(reqID,
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
