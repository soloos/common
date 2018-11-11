package types

import (
	"io"
)

func (p *Connection) ReadAcquire() {
	p.readMutex.Lock()
}

func (p *Connection) ReadRelease() {
	p.readMutex.Unlock()
}

func (p *Connection) AfterReadHeaderError() {
	p.ContinueReadSig.Done()
	p.ReadRelease()
}

func (p *Connection) AfterReadHeaderSuccess() error {
	if p.LastReadLimit == 0 {
		p.ContinueReadSig.Done()
		p.ReadRelease()
	}
	return nil
}

func (p *Connection) innerAfterReadHeader(maxMessageLength, contentLen uint32, netVersion byte) error {
	// prepare & check request
	p.LastReadLimit = contentLen
	if netVersion != SNetVersion {
		return ErrWrongVersion
	}

	if p.LastReadLimit > maxMessageLength {
		return ErrMessageTooLong
	}

	return nil
}

func (p *Connection) ReadRequestHeader(maxMessageLength uint32, header *RequestHeader) error {
	p.ReadAcquire()
	p.ContinueReadSig.Add(1)

	var (
		offset, n int
		err       error
	)

	for offset = 0; offset < len(header); offset += n {
		n, err = p.NetConn.Read(header[offset:len(header)])
		if err != nil && err != io.EOF {
			p.ReadRelease()
			p.ContinueReadSig.Done()
			return err
		}
	}

	return p.innerAfterReadHeader(maxMessageLength, header.ContentLen(), header.Version())
}

func (p *Connection) ReadResponseHeader(maxMessageLength uint32, header *ResponseHeader) error {
	p.ReadAcquire()
	p.ContinueReadSig.Add(1)

	var (
		offset, n int
		err       error
	)

	for offset = 0; offset < len(header); offset += n {
		n, err = p.NetConn.Read(header[offset:len(header)])
		if err != nil && err != io.EOF {
			p.ReadRelease()
			p.ContinueReadSig.Done()
			return err
		}
	}

	return p.innerAfterReadHeader(maxMessageLength, header.ContentLen(), header.Version())
}

func (p *Connection) Read(b []byte) (int, error) {
	if p.LastReadLimit == 0 {
		return 0, io.EOF
	}

	var (
		n   int
		err error
	)

	if p.LastReadLimit > uint32(len(b)) {
		n, err = p.NetConn.Read(b)
	} else {
		n, err = p.NetConn.Read(b[:int(p.LastReadLimit)])
	}

	if err != nil {
		p.ContinueReadSig.Done()
		p.ReadRelease()
		return n, err
	}

	p.LastReadLimit -= uint32(n)
	if p.LastReadLimit == 0 {
		p.ContinueReadSig.Done()
		p.ReadRelease()
	}

	return n, err
}

func (p *Connection) ReadAll(b []byte) error {
	var (
		off, n int
		err    error
	)

	for off = 0; off < len(b); off += n {
		n, err = p.Read(b[off:])
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
	}
	return nil
}
