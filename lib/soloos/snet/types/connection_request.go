package types

func (p *Connection) WriteRequestHeader(reqID uint64, serviceID string, bodySize, reqParamSize uint32) error {
	p.WriteAcquire()

	var (
		header RequestHeader
		off, n int
		err    error
	)
	header.SetID(reqID)
	header.SetVersion(SNetVersion)
	header.SetBodySize(bodySize)
	header.SetParamSize(reqParamSize)
	header.SetServiceID(serviceID)
	for off = 0; off < len(header); off += n {
		n, err = p.NetConn.Write(header[:])
		if err != nil {
			p.WriteRelease()
			return err
		}
	}

	return p.afterWriteHeader(bodySize)
}
