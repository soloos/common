package snet

func (p *NetQuery) writeSNetReqHeaderWithLock(reqID uint64, url string, bodySize, reqParamSize uint32) error {
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
			return err
		}
	}

	for off = 0; off < len(reqHeader.SNetReqHeaderBodyBs); off += n {
		n, err = p.Conn.NetConn.Write(reqHeader.SNetReqHeaderBodyBs)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *NetQuery) Request(req *SNetReq) error {
	p.Conn.WriteAcquire()

	var (
		err error
	)
	// post data
	err = req.NetQuery.writeSNetReqHeaderWithLock(req.ReqID,
		req.Url,
		uint32(len(req.Param)+(req.OffheapBody.BodySize())),
		uint32(len(req.Param)))
	if err != nil {
		goto POST_DATA_DONE
	}

	err = req.NetQuery.WriteAll(req.Param)
	if err != nil {
		goto POST_DATA_DONE
	}

	err = req.OffheapBody.Copy(&req.NetQuery)
	if err != nil {
		goto POST_DATA_DONE
	}

POST_DATA_DONE:
	if err != nil {
		err = req.NetQuery.ConnClose(err)
	}
	return err
}
