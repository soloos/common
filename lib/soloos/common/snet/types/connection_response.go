package types

func (p *Connection) WriteResponseHeader(reqID uint64, bodySize, reqParamSize uint32) error {
	p.WriteAcquire()

	var (
		header ResponseHeader
		off, n int
		err    error
	)
	header.SetID(reqID)
	header.SetVersion(SNetVersion)
	header.SetBodySize(bodySize)
	header.SetParamSize(reqParamSize)
	for off = 0; off < len(header); off += n {
		n, err = p.NetConn.Write(header[:])
		if err != nil {
			p.WriteRelease()
			return err
		}
	}

	return p.afterWriteHeader(bodySize)
}

func (p *Connection) SimpleResponse(reqID uint64, respBody []byte) error {
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

func (p *Connection) ResponseHeaderParam(reqID uint64, param []byte, offheapBodySize int) error {
	var err error
	err = p.WriteResponseHeader(reqID,
		uint32(len(param)+offheapBodySize),
		uint32(len(param)))
	if err != nil {
		return err
	}

	err = p.WriteAll(param)
	if err != nil {
		goto POST_DATA_DONE
	}

POST_DATA_DONE:
	return err
}

func (p *Connection) Response(reqID uint64, param []byte, offheapBody []byte) error {
	var err error
	err = p.WriteResponseHeader(reqID,
		uint32(len(param)+len(offheapBody)),
		uint32(len(param)))
	if err != nil {
		return err
	}

	err = p.WriteAll(param)
	if err != nil {
		goto POST_DATA_DONE
	}

	err = p.WriteAll(offheapBody)
	if err != nil {
		goto POST_DATA_DONE
	}

POST_DATA_DONE:
	return err
}
